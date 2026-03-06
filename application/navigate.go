package application

import "github.com/tahrioui/code-walkthrough/domain"

type NavigateUseCase struct {
	nav         *domain.Navigator
	walkthrough domain.Walkthrough
}

func NewNavigateUseCase(w domain.Walkthrough) *NavigateUseCase {
	return &NavigateUseCase{
		nav:         domain.NewNavigator(w),
		walkthrough: w,
	}
}

func (uc *NavigateUseCase) Current() (domain.Step, error) {
	return uc.nav.Current()
}

func (uc *NavigateUseCase) StepForward() (domain.Step, error) {
	return uc.nav.Next()
}

func (uc *NavigateUseCase) StepBackward() (domain.Step, error) {
	return uc.nav.Prev()
}

func (uc *NavigateUseCase) JumpTo(id domain.StepID) (domain.Step, error) {
	return uc.nav.JumpTo(id)
}

func (uc *NavigateUseCase) JumpToSection(id domain.SectionID) (domain.Step, error) {
	return uc.nav.JumpToSection(id)
}

func (uc *NavigateUseCase) CurrentSection() domain.Section {
	return uc.nav.CurrentSection()
}

func (uc *NavigateUseCase) CurrentIndex() int {
	return uc.nav.CurrentIndex()
}

func (uc *NavigateUseCase) TotalSteps() int {
	return uc.nav.TotalSteps()
}

func (uc *NavigateUseCase) ViewTOC() []domain.Section {
	return uc.walkthrough.Sections
}
