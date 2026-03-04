package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoSteps         = errors.New("walkthrough has no steps")
	ErrAtEnd           = errors.New("already at last step")
	ErrAtStart         = errors.New("already at first step")
	ErrStepNotFound    = errors.New("step not found")
	ErrSectionNotFound = errors.New("section not found")
)

type Navigator struct {
	walkthrough Walkthrough
	flatSteps   []Step
	sectionMap  map[int]SectionID // flatSteps index -> section ID
	index       int
}

func NewNavigator(w Walkthrough) *Navigator {
	var flat []Step
	sectionMap := make(map[int]SectionID)
	for _, sec := range w.Sections {
		for _, step := range sec.Steps {
			sectionMap[len(flat)] = sec.ID
			flat = append(flat, step)
		}
	}
	return &Navigator{
		walkthrough: w,
		flatSteps:   flat,
		sectionMap:  sectionMap,
		index:       0,
	}
}

func (n *Navigator) Current() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	return n.flatSteps[n.index], nil
}

func (n *Navigator) CurrentIndex() int {
	return n.index
}

func (n *Navigator) Next() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	if n.index >= len(n.flatSteps)-1 {
		return Step{}, ErrAtEnd
	}
	n.index++
	return n.flatSteps[n.index], nil
}

func (n *Navigator) Prev() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	if n.index <= 0 {
		return Step{}, ErrAtStart
	}
	n.index--
	return n.flatSteps[n.index], nil
}

func (n *Navigator) JumpTo(id StepID) (Step, error) {
	for i, step := range n.flatSteps {
		if step.ID == id {
			n.index = i
			return step, nil
		}
	}
	return Step{}, fmt.Errorf("%w: %s", ErrStepNotFound, id)
}

func (n *Navigator) JumpToSection(id SectionID) (Step, error) {
	for i, secID := range n.sectionMap {
		if secID == id {
			n.index = i
			return n.flatSteps[i], nil
		}
	}
	return Step{}, fmt.Errorf("%w: %s", ErrSectionNotFound, id)
}

func (n *Navigator) CurrentSection() Section {
	secID := n.sectionMap[n.index]
	for _, sec := range n.walkthrough.Sections {
		if sec.ID == secID {
			return sec
		}
	}
	return Section{}
}

func (n *Navigator) TotalSteps() int {
	return len(n.flatSteps)
}
