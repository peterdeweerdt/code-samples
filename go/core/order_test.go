package core_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"core"
	"pos"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderForPager(t *testing.T) {
	const (
		siteID      = core.KountaID(123)
		pagerNumber = 765
	)

	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	order, err := app.CreateOrderForPager(siteID, pagerNumber)
	assert.NoError(t, err)

	order, err = app.FindOrderByID(order.ID)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", pagerNumber), order.PagerNumber)
	assert.Equal(t, siteID, order.SiteID)
}

func TestLinkOrderWithTable(t *testing.T) {
	const (
		siteID      = core.KountaID(123)
		pagerNumber = 765
		tableName   = "45"
	)

	var app core.AppContext
	defer testServer(&app)()

	order := &core.Order{SiteID: siteID, PagerNumber: fmt.Sprintf("%d", pagerNumber)}
	err := app.DB.InsertOrder(order)
	assert.NoError(t, err)
	assert.NotZero(t, order.ID)

	err = app.LinkOrderWithTable(siteID, pagerNumber, tableName)
	assert.NoError(t, err)
	order, err = app.FindOrderByID(order.ID)
	assert.NoError(t, err)
	assert.Equal(t, tableName, order.TableName)
}

func TestFindPayableOrdersByTableNameWillFindOnHoldOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	posOrder := pos.NewMockKountaOrder()

	_, err := app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrders, err := app.FindPayableOrdersByTableName(expectedOrder.SiteID, expectedOrder.TableName)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrders[0])
	assertOrdersEqual(t, expectedOrder, &actualOrders[0])
}

func TestFindPayableOrdersByTableNameWillFindPendingOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	expectedOrder.Status = core.OrderStatusPending
	posOrder := pos.NewMockKountaOrder()
	posOrder.Status = core.OrderStatusPending

	_, err := app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrders, err := app.FindPayableOrdersByTableName(expectedOrder.SiteID, expectedOrder.TableName)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrders[0])
	assertOrdersEqual(t, expectedOrder, &actualOrders[0])
}

func TestFindPayableOrdersByPagerIdWillFindOnHoldOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	pagerNumber, err := strconv.ParseInt(expectedOrder.PagerNumber, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	posOrder := pos.NewMockKountaOrder()

	_, err = app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrders, err := app.FindPayableOrdersByPagerID(expectedOrder.SiteID, pagerNumber)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrders[0])
	assertOrdersEqual(t, expectedOrder, &actualOrders[0])
}

func TestFindPayableOrdersByPagerIdWillFindPendingOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	expectedOrder.Status = core.OrderStatusPending
	pagerNumber, err := strconv.ParseInt(expectedOrder.PagerNumber, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	posOrder := pos.NewMockKountaOrder()
	posOrder.Status = core.OrderStatusPending

	_, err = app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrders, err := app.FindPayableOrdersByPagerID(expectedOrder.SiteID, pagerNumber)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrders[0])
	assertOrdersEqual(t, expectedOrder, &actualOrders[0])
}

func TestFindOrderByID(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	posOrder := pos.NewMockKountaOrder()

	_, err := app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrder, err := app.FindOrderByID(expectedOrder.ID) // should be only one order in DB
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrder)
	assertOrdersEqual(t, expectedOrder, actualOrder)
}

func TestFindOrderByPosID(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	posOrder := pos.NewMockKountaOrder()

	_, err := app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// act
	actualOrder, err := app.FindOrderByPosID(posOrder.GetPosID())
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, actualOrder)
	assertOrdersEqual(t, expectedOrder, actualOrder)
}

func TestCreateOrderFromPosOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()
	posOrder := pos.NewMockKountaOrder()

	// act
	actualOrder, err := app.CreateOrUpdateOrderFromKounta(posOrder)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assertOrdersEqual(t, expectedOrder, actualOrder)
}

func TestUpdateOrderFromPosOrder(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()

	app.TestInsertMenu(t)

	expectedOrder := core.NewExpectedOrder()

	expectedOrder.Status = "ACCEPTED"
	expectedOrder.Lines[0].AddedModifiers[0].LineID = 3 // line id will update
	expectedOrder.Lines[1].RemovedModifiers[0].LineID = 4

	initialPosOrder := pos.NewMockKountaOrder()

	// act
	_, err := app.CreateOrUpdateOrderFromKounta(initialPosOrder)
	if err != nil {
		t.Fatal(err)
	}

	updatedPosOrder := pos.NewMockKountaOrder()
	updatedPosOrder.Status = "ACCEPTED"
	actualOrder, err := app.CreateOrUpdateOrderFromKounta(updatedPosOrder)

	// assert
	assertOrdersEqual(t, expectedOrder, actualOrder)
}

func TestRejectedOrderClearsPager(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	initialPosOrder := pos.NewMockKountaOrder()
	initialPosOrder.Status = core.OrderStatusSubmitted
	app.CreateOrUpdateOrderFromKounta(initialPosOrder)

	//act
	updatedPosOrder := pos.NewMockKountaOrder()
	updatedPosOrder.Status = core.OrderStatusRejected
	actualOrder, err := app.CreateOrUpdateOrderFromKounta(updatedPosOrder)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "", actualOrder.PagerNumber)
}

func TestCompletedOrderClearsPager(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	initialPosOrder := pos.NewMockKountaOrder()
	initialPosOrder.Status = core.OrderStatusPending
	app.CreateOrUpdateOrderFromKounta(initialPosOrder)

	//act
	updatedPosOrder := pos.NewMockKountaOrder()
	updatedPosOrder.Status = core.OrderStatusComplete
	actualOrder, _ := app.CreateOrUpdateOrderFromKounta(updatedPosOrder)

	// assert
	assert.Equal(t, "", actualOrder.PagerNumber)
}

func TestDeletedOrderClearsPager(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	initialPosOrder := pos.NewMockKountaOrder()
	initialPosOrder.Status = core.OrderStatusAccepted
	app.CreateOrUpdateOrderFromKounta(initialPosOrder)

	//act
	updatedPosOrder := pos.NewMockKountaOrder()
	updatedPosOrder.Status = core.OrderStatusDeleted
	actualOrder, _ := app.CreateOrUpdateOrderFromKounta(updatedPosOrder)

	// assert
	assert.Equal(t, "", actualOrder.PagerNumber)
}

func TestUpdatePickupDetailsUpdatesDatabase(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	initialPosOrder := pos.NewMockKountaOrder()
	initialPosOrder.Status = core.OrderStatusSubmitted
	order, _ := app.CreateOrUpdateOrderFromKounta(initialPosOrder)

	if mockKounta, ok := app.Kounta.(*pos.MockKounta); ok {
		mockKounta.Orders = append(mockKounta.Orders, *initialPosOrder)
	}

	pickupDetails := core.PickupDetails{}
	byteBody := []byte(`{"customer_name":"John Doe","phone_number":"4041234567","pickup_time":"2017-04-21T18:05:37+02:00"}`)
	json.Unmarshal(byteBody, &pickupDetails)

	// act
	err := app.UpdatePickupDetails(order.ID, pickupDetails)

	// assert
	assert.NoError(t, err)
	updatedOrder, err := app.DB.GetOrderByDatabaseID(order.ID)
	expectedPickupTime := pickupDetails.PickupTime
	assert.True(t, updatedOrder.PickupTime.Equal(*expectedPickupTime))
}

func TestUpdatePickupDetailsUpdatesKountaNotes(t *testing.T) {
	// arrange
	var app core.AppContext
	defer testServer(&app)()
	app.TestInsertMenu(t)

	initialPosOrder := pos.NewMockKountaOrder()
	initialPosOrder.Status = core.OrderStatusSubmitted
	order, _ := app.CreateOrUpdateOrderFromKounta(initialPosOrder)

	if mockKounta, ok := app.Kounta.(*pos.MockKounta); ok {
		mockKounta.Orders = append(mockKounta.Orders, *initialPosOrder)
	}

	pickupDetails := core.PickupDetails{}
	byteBody := []byte(`{"customer_name":"John Doe","phone_number":"4041234567","pickup_time":"2017-04-21T18:05:37+02:00"}`)
	json.Unmarshal(byteBody, &pickupDetails)

	// act
	err := app.UpdatePickupDetails(order.ID, pickupDetails)

	// assert
	assert.NoError(t, err)
	if mockKounta, ok := app.Kounta.(*pos.MockKounta); ok {
		assert.Equal(t, "TO-GO APP - PAID\n\n6:05 pm\n\nJohn Doe\n\n4041234567", mockKounta.Notes)
	} else {
		t.Fatal("Kounta struct used in testing is not of type *MockKounta.")
	}
}

func TestAppendUniqueOrdersWhenAllNewOrders(t *testing.T) {
	// arrange
	var app core.AppContext

	order1ID := core.DatabaseID(1)
	order2ID := core.DatabaseID(2)
	order3ID := core.DatabaseID(3)
	order4ID := core.DatabaseID(4)

	order1 := core.Order{ID: order1ID}
	order2 := core.Order{ID: order2ID}
	order3 := core.Order{ID: order3ID}
	order4 := core.Order{ID: order4ID}
	initial := []core.Order{order1, order2}
	newOrders := []core.Order{order3, order4}

	// act
	updated := app.AppendUniqueOrders(initial, newOrders)

	// assert
	assert.Equal(t, 4, len(updated))
	assert.Equal(t, order1ID, updated[0].ID)
	assert.Equal(t, order2ID, updated[1].ID)
	assert.Equal(t, order3ID, updated[2].ID)
	assert.Equal(t, order4ID, updated[3].ID)
}

func TestAppendUniqueOrdersWhenSomeNewOrders(t *testing.T) {
	// arrange
	var app core.AppContext

	order1ID := core.DatabaseID(1)
	order2ID := core.DatabaseID(2)
	order3ID := core.DatabaseID(3)

	order1 := core.Order{ID: order1ID}
	order2 := core.Order{ID: order2ID}
	order3 := core.Order{ID: order3ID}
	initial := []core.Order{order1, order2}
	newOrders := []core.Order{order1, order3}

	// act
	updated := app.AppendUniqueOrders(initial, newOrders)

	// assert
	assert.Equal(t, 3, len(updated))
	assert.Equal(t, order1ID, updated[0].ID)
	assert.Equal(t, order2ID, updated[1].ID)
	assert.Equal(t, order3ID, updated[2].ID)
}

func TestAppendUniqueOrdersWhenNoNewOrders(t *testing.T) {
	// arrange
	var app core.AppContext

	order1ID := core.DatabaseID(1)
	order2ID := core.DatabaseID(2)

	order1 := core.Order{ID: order1ID}
	order2 := core.Order{ID: order2ID}
	initial := []core.Order{order1, order2}

	// act
	updated := app.AppendUniqueOrders(initial, initial)

	// assert
	assert.Equal(t, 2, len(updated))
	assert.Equal(t, order1ID, updated[0].ID)
	assert.Equal(t, order2ID, updated[1].ID)
}

func TestAppendUniqueOrdersWhenDuplicateNewOrders(t *testing.T) {
	// arrange
	var app core.AppContext

	order1ID := core.DatabaseID(1)
	order2ID := core.DatabaseID(2)
	order3ID := core.DatabaseID(3)

	order1 := core.Order{ID: order1ID}
	order2 := core.Order{ID: order2ID}
	order3 := core.Order{ID: order3ID}
	initial := []core.Order{order1, order2}
	newOrders := []core.Order{order3, order3}

	// act
	updated := app.AppendUniqueOrders(initial, newOrders)

	// assert
	assert.Equal(t, 3, len(updated))
	assert.Equal(t, order1ID, updated[0].ID)
	assert.Equal(t, order2ID, updated[1].ID)
	assert.Equal(t, order3ID, updated[2].ID)
}

// helpers

func assertOrdersEqual(t *testing.T, expected *core.Order, actual *core.Order) {
	if expected == nil || actual == nil {
		assert.Fail(t, "expected and actual orders cannot be nil")
	}

	assert.Equal(t, expected.TableName, actual.TableName, "table name should be equal")
	assert.Equal(t, expected.PosID, actual.PosID, "POS ID should be equal")
	assert.Equal(t, expected.Total, actual.Total, "total should be equal")
	assert.Equal(t, len(expected.Lines), len(actual.Lines), "Lines length should be equal")
	assert.Equal(t, expected.Lines[0].Notes, actual.Lines[0].Notes, "notes should be equal")
	assert.Equal(t, expected.Status, actual.Status, "statuses must be equal")

	// check modifiers on line 0
	assertModifiersEqual(t, expected.Lines[0].AddedModifiers, actual.Lines[0].AddedModifiers)
	assertModifiersEqual(t, expected.Lines[0].RemovedModifiers, actual.Lines[0].RemovedModifiers)
	// check modifiers on line 1
	assertModifiersEqual(t, expected.Lines[1].AddedModifiers, actual.Lines[1].AddedModifiers)
	assertModifiersEqual(t, expected.Lines[1].RemovedModifiers, actual.Lines[1].RemovedModifiers)
}

func assertModifiersEqual(t *testing.T, expected []core.Modifier, actual []core.Modifier) {
	if (expected == nil || actual == nil) && !(expected == nil && actual == nil) { // XOR nil check

		assert.Fail(t, "if one set of modifiers is nil, so must the other")
		return
	}

	if len(expected) != len(actual) {
		assert.Fail(t, "modifier lists must be same length")
		return
	}

	for i := 0; i < len(expected); i++ {
		assert.True(t, expected[i].Equals(actual[i]), "modifiers must be equal")
	}
}
