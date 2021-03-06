// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

type GraphicsContext interface {
	SetSize(width, height int, scale float64) error
	Update() error
}

type UserInterface interface {
	Run(width, height int, scale float64, title string, g GraphicsContext) error
	ScreenScale() float64
	SetScreenSize(width, height int) (bool, error)
	SetScreenScale(scale float64) (bool, error)
}

type RegularTermination struct {
}

func (*RegularTermination) Error() string {
	return "regular termination"
}
