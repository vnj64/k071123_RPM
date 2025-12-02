package permissions

import (
	"fmt"
	"k071123/pkg/enum"
)

type Permission string

// PermissionMeta метаданные пермишена.
// В будущем мб прикрутим сюда код ошибки или сообщение о причине пермишена.
// Пока что это просто закос под возможное масштабирование
type PermissionMeta struct {
	Description string
}

var (
	permissions     = enum.New[string, Permission]()
	permissionsData = make(map[Permission]PermissionMeta, 100)
)

func register(enumCase Permission, description ...string) Permission {
	if _, exists := permissions.From(enumCase.Value()); exists {
		panic(fmt.Sprintf("Permission %s is already registered", enumCase))
	}

	if len(description) > 0 {
		permissionsData[enumCase] = PermissionMeta{Description: description[0]}
	} else {
		permissionsData[enumCase] = PermissionMeta{Description: ""}
	}

	return permissions.Register(enumCase)
}

var (
	OnlySuperAdmin = register("any", "any")

	// User Service
	CreateUser    = register("user.user.create", "Create user")
	UpdateProfile = register("user.user.update", "Self user update")

	// Parking Service
	CreateParking = register("parking.parking.create", "Create Parking")

	CreateTariff = register("parking.tariff.create", "Create Tariff")

	StartSession = register("parking.session.start", "Start Parking Session")
	StopSession  = register("parking.session.stop", "Stop Parking Session")

	CreateUnit = register("parking.unit.create", "Create Unit")

	// Notification Service
	// Order Service
	SaveCard   = register("order.card.save", "Save Card")
	VerifyCard = register("order.card.verify", "Verify Card")
)

type RolePermissions map[Role][]Permission

var rolePermissions = RolePermissions{
	SuperAdmin: {
		CreateUser, UpdateProfile,
		CreateParking,
		CreateTariff,
		CreateUnit,
		StopSession,
	},
	ParkingAdmin: {
		UpdateProfile,
		CreateParking,
		CreateTariff,
		StopSession,
		CreateUnit,
	},
	Default: {
		UpdateProfile,
		StartSession,
		StopSession,
		SaveCard,
		VerifyCard,
	},
}

// PermissionsForRole Возвращает разрешения для конкретной роли
func PermissionsForRole(role Role) []Permission {
	return rolePermissions[role]
}

// Методы для работы с разрешениями

// Value Возвращает строковое значение разрешения
func (s Permission) Value() string {
	return string(s)
}

// All Получить все разрешения
func All() []string {
	return permissions.Keys()
}

// Enum Возвращает интерфейс для чтения enum
func Enum() enum.Reader[string, Permission] {
	return permissions
}

// Description Получить описание для разрешения
func (s Permission) Description() (string, bool) {
	permMeta, ok := permissionsData[s]
	if !ok {
		return "", false
	}
	return permMeta.Description, true
}

// FromString Получить разрешение из строки
func FromString(candidate string) (Permission, error) {
	permission, ok := permissions.From(candidate)
	if !ok {
		return "", fmt.Errorf("permission %s not found", candidate)
	}
	return permission, nil
}
