package checkpoint

import (
	"github.com/emicklei/go-restful"
	"k8s.io/kubernetes/pkg/kubelet/container"
	"net/http"
	"path"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	klog "k8s.io/klog/v2"
	kubepod "k8s.io/kubernetes/pkg/kubelet/pod"
)

type Manager interface {
	HandleCheckpointRequest(*restful.Request, *restful.Response)
	FindCheckpointForPod(*v1.Pod) (Checkpoint, bool)
	// TriggerPodMigration(*v1.Pod) (Result, error)
}

type Checkpoint interface {
	WaitUntilFinished()
	Options() *container.CheckpointPodOptions
}

func NewManager(kubeClient clientset.Interface, podManager kubepod.Manager, checkpointFn checkpointFunc, rootPath string) Manager {
	return &manager{
		checkpointPath:     path.Join(rootPath, "checkpoint"),
		kubeClient:         kubeClient,
		podManager:         podManager,
		checkpointFc:       checkpointFn,
		checkpoints:        make(map[types.UID]*checkpoint),
	}
}

type checkpointFunc func(*v1.Pod)

type manager struct {
	checkpointPath string

	kubeClient     clientset.Interface
	podManager     kubepod.Manager
	checkpointFc   checkpointFunc

	checkpoints    map[types.UID]*checkpoint
}

var _ Manager = &manager{}

type checkpoint struct {
	path       string
	containers []string
	unblock    chan struct{}
	done       chan struct{}
}

type Result struct {
	Containers map[string]ResultContainer
}

type ResultContainer struct {
	CheckpointPath string
}

var _ Checkpoint = &checkpoint{}

func (m *manager) HandleCheckpointRequest(req *restful.Request, res *restful.Response) {
	params := getCheckpointRequestParams(req)
	klog.V(2).Infof("POST Checkpoint - %v %v", params.podUID, params.containerNames)

	var pod *v1.Pod
	var ok bool

	if pod, ok = m.podManager.GetPodByUID(types.UID(params.podUID)); !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if pod.Status.Phase != v1.PodRunning {
		res.WriteHeader(http.StatusConflict)
		return
	}

	ck := m.newCheckpoint(pod)
	ck.containers = params.containerNames

	klog.V(2).Infof("Starting migration of Pod %v", pod.Name)
	m.checkpointFc(pod)

	<-ck.done
	r := Result{Containers: map[string]ResultContainer{}}
	for _, c := range ck.containers {
		r.Containers[c] = ResultContainer{CheckpointPath: path.Join(ck.path, c)}
	}
	if err := res.WriteAsJson(r); err != nil {
		klog.Error("failed to encode migration result.", err)
	}
	res.WriteHeader(http.StatusOK)
	ck.unblock <- struct{}{}
}

func (m *manager) FindCheckpointForPod(pod *v1.Pod) (Checkpoint, bool) {
	ck, ok := m.checkpoints[pod.UID]
	return ck, ok
}

func (m *manager) newCheckpoint(pod *v1.Pod) *checkpoint {
	ck := &checkpoint{
		path:    path.Join(m.checkpointPath, pod.Name),
		unblock: make(chan struct{}),
		done:    make(chan struct{}),
	}
	m.checkpoints[pod.GetUID()] = ck
	return ck
}

func (m *manager) removeCheckpoint(pod *v1.Pod) {
	mig, ok := m.checkpoints[pod.GetUID()]
	if !ok {
		return
	}
	mig.done <- struct{}{}
	delete(m.checkpoints, pod.GetUID())
}

func (ck *checkpoint) Options() *container.CheckpointPodOptions {
	return &container.CheckpointPodOptions{
		KeepRunning:    true,
		CheckpointsDir: ck.path,
		Unblock:        ck.unblock,
		Done:           ck.done,
		Containers:     ck.containers,
	}
}

func (ck *checkpoint) WaitUntilFinished() {
	<-ck.unblock
}

type CheckpointRequestParams struct {
	podUID         string
	containerNames []string
}

func getCheckpointRequestParams(req *restful.Request) CheckpointRequestParams {
	return CheckpointRequestParams{
		podUID:         req.PathParameter("podUID"),
		containerNames: strings.Split(req.QueryParameter("containers"), ","),
	}
}
