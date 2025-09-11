package interfaces

import (
	"reflect"

	"github.com/control-monkey/terraform-provider-cm/internal/helpers"
	"github.com/hashicorp/go-set/v2"
)

type OperationType string

const (
	CreateOperation OperationType = "create"
	UpdateOperation OperationType = "update"
	DeleteOperation OperationType = "delete"
)

type MergedEntities[T MergeModel] struct {
	EntitiesToCreate set.Collection[T]
	EntitiesToUpdate set.Collection[T]
	EntitiesToDelete set.Collection[T]
}

type MergeModel interface {
	set.Hasher[string]
	GetBlockIdentifier() string
}

func GetIdentifiers[T MergeModel](ms []T) []string {
	knownEntities := filterOutUnknowns(ms)
	retVal := mapToIdentifiers(knownEntities)

	return retVal
}

func filterOutUnknowns[T MergeModel](ms []T) []T {
	f := func(m T) bool {
		return (m).GetBlockIdentifier() != ""
	}

	retVal := helpers.Filter(ms, f)
	return retVal

}

func mapToIdentifiers[T MergeModel](ms []T) []string {
	mapToIdentifier := func(m T) string {
		return (m).GetBlockIdentifier()
	}

	return helpers.Map(ms, mapToIdentifier)
}

func MergeEntities[T MergeModel](plan []T, state []T) MergedEntities[T] {
	planEntities := set.HashSetFrom[T, string](plan)
	stateEntities := set.HashSetFrom[T, string](state)

	entitiesToCreate := planEntities.Difference(stateEntities)
	entitiesToDelete := stateEntities.Difference(planEntities)

	idToModelToDelete := make(map[string]T, 0)
	entitiesToRemoveFromCreate := make([]T, 0)
	entitiesToRemoveFromDelete := make([]T, 0)
	entitiesToUpdate := set.NewHashSet[T, string](0)

	for _, p := range entitiesToDelete.Slice() {
		idToModelToDelete[p.GetBlockIdentifier()] = p
	}
	for _, p := range entitiesToCreate.Slice() {
		if entityToDelete := idToModelToDelete[p.GetBlockIdentifier()]; reflect.ValueOf(entityToDelete).IsNil() == false { // id in both add & delete
			entitiesToRemoveFromCreate = append(entitiesToRemoveFromCreate, p)
			entitiesToRemoveFromDelete = append(entitiesToRemoveFromDelete, entityToDelete)
			entitiesToUpdate.Insert(p)
		}
	}

	entitiesToCreate.RemoveSlice(entitiesToRemoveFromCreate)
	entitiesToDelete.RemoveSlice(entitiesToRemoveFromDelete)

	retVal := MergedEntities[T]{
		EntitiesToCreate: entitiesToCreate,
		EntitiesToUpdate: entitiesToUpdate,
		EntitiesToDelete: entitiesToDelete,
	}

	return retVal
}
