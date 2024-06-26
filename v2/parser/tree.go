package parser

import (
	"strings"

	objs "github.com/SakoDroid/telego/v2/objects"
)

// TreeNode is a special tree element containing handlers.
type TreeNode struct {
	father      *TreeNode
	right, left *TreeNode
	data        *handler
}

// handlerTree is a special binary tree for storing handlers. Right node hase a value that does not match the it's father regex and the left node matches it's father regex.
type handlerTree struct {
	root *TreeNode
}

// AddHandler adds a new handler to the tree.
func (tr *handlerTree) AddHandler(hdl *handler) {
	tn := TreeNode{data: hdl}
	tr.addNode(&tn)
}

// GetHandler gets the proper handler for the given text.
func (tr *handlerTree) GetHandler(msg *objs.Message) *handler {
	msgText := msg.Text
	if msg.Caption != "" {
		msgText = msg.Caption
	}
	tn := tr.findTheNodeRegex(msgText, msg.Chat.Type)
	if tn != nil {
		return tn.data
	}
	return nil
}

func (tr *handlerTree) findTheNodeRegex(text, chatType string) *TreeNode {
	node := tr.root
	for {
		if node == nil {
			break
		}
		if node.data.regex.Match([]byte(text)) {
			if node.left != nil {
				node = node.left
			} else {
				break
			}
		} else {
			if node.right != nil {
				node = node.right
			} else {
				node = node.father
				break
			}
		}
	}
	return tr.checkForChatTypes(node, chatType, text)
}

func (tr *handlerTree) checkForChatTypes(currentNode *TreeNode, chatType, text string) *TreeNode {
	for {
		if currentNode == nil {
			break
		}
		if (strings.Contains(currentNode.data.chatType, chatType) || strings.Contains(currentNode.data.chatType, "all")) && currentNode.data.regex.Match([]byte(text)) {
			break
		} else {
			currentNode = currentNode.father
		}
	}
	return currentNode
}

// Finds the perfect location for this handler.
func (tr *handlerTree) addNode(tn *TreeNode) {
	var fatherNode *TreeNode
	node := tr.root
	dir := 0
	for {
		//this is the spot
		if node == nil {
			if fatherNode != nil {
				tn.father = fatherNode
				if dir == 0 {
					fatherNode.left = tn
				} else {
					fatherNode.right = tn
				}
			} else {
				//Root node
				tr.root = tn
			}
			break
		} else {
			if tr.checkRegex(node, tn) {
				dir = 0
				fatherNode = node
				node = fatherNode.left
			} else {
				dir = 1
				fatherNode = node
				node = fatherNode.right
			}
		}
	}
}

func (tr *handlerTree) checkRegex(father, child *TreeNode) bool {
	return father.data.regex.Match(
		[]byte(child.data.regex.String()),
	)
}
