package scanner

import (
	"container/list"
	"fmt"
	"strings"
)

type TreeNode struct {
	Value  Token
	parent *TreeNode
	items  *list.List
}

// NewTreeElement Creates a new TreeElement.
func NewTreeNode(value Token) *TreeNode {
	return &TreeNode{value, nil, list.New()}
}

// Parent Returns the current element parent
func (this *TreeNode) Parent() *TreeNode {
	return this.parent
}

// setParent sets the current nodes parent value.
// Warning: does not add the node as a child
func (this *TreeNode) setParent(element *TreeNode) {
	if this.parent != nil {
		this.parent.RemoveChild(this)
		this.parent = nil
	}
	this.parent = element
}

func (this *TreeNode) LastElement() *TreeNode {
	if this.items.Len() == 0 {
		return nil
	}
	r1, ok := this.items.Back().Value.(*TreeNode)
	if !ok {
		panic("Expecting TreeElement as child node.")
	}
	return r1
}

func (this *TreeNode) Last() Token {
	last := this.LastElement()
	if last != nil {
		return last.Value
	}
	return nil
}

func (this *TreeNode) FirstChild() *list.Element {
	return this.items.Front()
}

// Add adds a TreeElement to the end of the children items of the current node.
func (this *TreeNode) AddElement(element *TreeNode) *TreeNode {
	element.setParent(this)
	this.items.PushBack(element)
	return element
}

// Add adds a value to the end of the children items of the current node.
func (this *TreeNode) Add(value Token) *TreeNode {
	element := NewTreeNode(value)
	return this.AddElement(element)
}

// Push, removes the current element from its current parent, place the new value
// in its place and add the current element to the new element. there by pushing the current
// element down the hierachy.
// Example:
// tree:  A(B)
// B.Push(C)
// tree:  A(C(B))
func (this *TreeNode) PushElement(element *TreeNode) *TreeNode {
	parent := this.Parent()
	if parent != nil {
		//replace the current node with the new node
		listelement := parent.findChildElement(this)
		listelement.Value = element
		element.setParent(parent)
		this.parent = nil
	}
	//add the current node to the new node
	element.AddElement(this)
	return element
}

func (this *TreeNode) Push(value Token) *TreeNode {
	return this.PushElement(NewTreeNode(value))
}

//FindChildElement Finds a child element in the current nodes children
func (this *TreeNode) findChildElement(element *TreeNode) *list.Element {
	for e := this.items.Front(); e != nil; e = e.Next() {
		if e.Value == element {
			return e
		}
	}
	return nil
}

//FindChild Finds a child element in the current nodes children
func (this *TreeNode) FindChildElement(element *TreeNode) *TreeNode {
	listelement := this.findChildElement(element)
	if listelement == nil {
		return nil
	}
	r1, ok := listelement.Value.(*TreeNode)
	if ok {
		return r1
	}
	panic("Wrong Type, expecting element type to be of type '*TreeElement'")
}

func (this *TreeNode) RemoveChild(element *TreeNode) {
	listelement := this.findChildElement(element)
	if listelement == nil {
		panic("Element not found.")
	}
	this.items.Remove(listelement)
}

//FindChildValue Finds a chile value in the current nodes children.
func (this *TreeNode) FindChild(value Token) *TreeNode {
	for e := this.items.Front(); e != nil; e = e.Next() {
		telement, ok := e.Value.(*TreeNode)
		if ok && telement.Value == value {
			return telement
		}
	}
	return nil
}

func (this *TreeNode) StringContent() string {
	lines := make([]string, this.items.Len())
	k := 0
	for item := this.items.Front(); item != nil; item = item.Next() {
		lines[k] = item.Value.(*TreeNode).String()
		k++
	}
	if this.Value.Error() != nil {
		return fmt.Sprintf("[ERROR: %s]", this.Value.Error())
	} else if len(lines) > 0 {
		return fmt.Sprintf("%s", strings.Join(lines, ","))
	} else {
		return ""
	}
}

func (this *TreeNode) String() string {
	if this.StringContent() == "" {
		return this.Value.String()
	}
	return fmt.Sprintf("[%s:%s]", this.Value.String(), this.StringContent())
}
