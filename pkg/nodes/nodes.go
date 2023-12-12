package nodes

import (
	"context"
	"encoding/json"
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type NodesInterface interface {
	AddNodesToRefresh() []Node
	CordonNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string) error
	DarinNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string) error
}

type Node struct {
	NodeName          string `json:"nodeName"`
	NodeVersion       string `json:"nodeVersion"`
	NodeArch          string `json:"nodeArch"`
	CreationTimestamp string `json:"creationTimestamp"`
}

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value bool   `json:"value"`
}

func enqueueNode(queue []Node, n Node) []Node {
	queue = append(queue, n)
	slog.Debug("Enqueued:" + n.NodeName)
	return queue
}

func dequeueNode(queue []Node) []Node {
	n := queue[0]
	slog.Debug("Dequeued:" + n.NodeName)
	return queue[1:]
}

func cordonNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string) error {

	nodePatch := []patchStringValue{
		{Op: "replace", Path: "/spec/unschedulable", Value: true},
	}
	payload, err := json.Marshal(nodePatch)
	if err != nil {
		return err
	}

	patchOptions := metav1.PatchOptions{}
	_, err = clientset.CoreV1().Nodes().Patch(ctx, nodeName, types.JSONPatchType, payload, patchOptions)
	return err
}

func addNodesToRefresh(ctx context.Context, clientset *kubernetes.Clientset, globalNodes []Node) []Node {

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		slog.Error("Error getting nodes", err)
	}

	for _, node := range nodes.Items {
		kubeletVersion := node.Status.NodeInfo.KubeletVersion
		nodeArch := node.Status.NodeInfo.Architecture
		nodeName := node.ObjectMeta.Name
		creationTimestamp := node.ObjectMeta.CreationTimestamp.String()

		node := Node{
			NodeName:          nodeName,
			NodeVersion:       kubeletVersion,
			NodeArch:          nodeArch,
			CreationTimestamp: creationTimestamp,
		}

		globalNodes = enqueueNode(globalNodes, node)
	}

	return globalNodes
}
