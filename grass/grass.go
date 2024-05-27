package grass

import "ecosim/entity"

type Grass interface {
	entity.Animal
	Update()
}
