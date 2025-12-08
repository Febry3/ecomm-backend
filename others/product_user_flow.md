# Product & Inventory User Flow

This document outlines the user flows for the Product and Inventory system based on the current database schema.

## 1. Seller: Product Creation & Inventory Management

### Flow Description
1.  **Create Product**: Seller defines the main product details (Title, Description, Category).
2.  **Create Variants**: Seller adds variations (e.g., Size S, Color Red).
    *   *System Action*: Automatically initializes `ProductVariantStock` with 0 stock.
3.  **Initial Stocking**: Seller adds initial stock to the variants.
    *   *System Action*: Updates `ProductVariantStock` and records `InventoryLedger`.

### Diagram
```mermaid
sequenceDiagram
    actor Seller
    participant API
    participant ProductRepo
    participant InventoryRepo
    participant DB

    Note over Seller, DB: Step 1: Create Product
    Seller->>API: POST /products (Title, Desc, Category)
    API->>ProductRepo: CreateProduct()
    ProductRepo->>DB: Insert into `products`
    DB-->>API: Product Created (ID)

    Note over Seller, DB: Step 2: Add Variants
    Seller->>API: POST /products/:id/variants (SKU, Name, Price)
    API->>ProductRepo: CreateVariant()
    ProductRepo->>DB: Insert into `product_variants`
    ProductRepo->>DB: Insert into `product_variant_stocks` (Default 0)
    DB-->>API: Variant Created (ID)

    Note over Seller, DB: Step 3: Add Stock
    Seller->>API: POST /inventory/add (VariantID, Qty, Reason)
    API->>InventoryRepo: UpdateStock(VariantID, +Qty)
    InventoryRepo->>DB: BEGIN TRANSACTION
    InventoryRepo->>DB: UPDATE `product_variant_stocks` SET current_stock += Qty
    InventoryRepo->>DB: INSERT INTO `inventory_ledgers` (Qty, Reason)
    InventoryRepo->>DB: COMMIT
    DB-->>Seller: Stock Updated
```

## 2. Buyer: Browsing & Purchasing

### Flow Description
1.  **View Product**: Buyer sees product details and available variants.
2.  **Check Availability**: System checks `CurrentStock` - `ReservedStock`.
3.  **Add to Cart**: Buyer adds item to cart.
4.  **Checkout (Reservation)**: When order is placed, stock is reserved.
    *   *System Action*: `ReservedStock` increases.
5.  **Payment Success (Deduction)**: When payment is confirmed, stock is permanently deducted.
    *   *System Action*: `CurrentStock` decreases, `ReservedStock` decreases.

### Diagram
```mermaid
sequenceDiagram
    actor Buyer
    participant API
    participant InventoryRepo
    participant OrderRepo
    participant DB

    Note over Buyer, DB: Step 1: View & Check Stock
    Buyer->>API: GET /products/:id
    API->>InventoryRepo: GetStock(VariantID)
    InventoryRepo->>DB: SELECT current_stock, reserved_stock
    DB-->>Buyer: Available Qty (Current - Reserved)

    Note over Buyer, DB: Step 2: Place Order (Reservation)
    Buyer->>API: POST /orders (VariantID, Qty)
    API->>InventoryRepo: ReserveStock(VariantID, Qty)
    InventoryRepo->>DB: UPDATE `product_variant_stocks` SET reserved_stock += Qty
    API->>OrderRepo: CreateOrder()
    DB-->>Buyer: Order Pending Payment

    Note over Buyer, DB: Step 3: Payment (Fulfillment)
    Buyer->>API: POST /payments/callback (Success)
    API->>InventoryRepo: DeductStock(VariantID, Qty)
    InventoryRepo->>DB: BEGIN TRANSACTION
    InventoryRepo->>DB: UPDATE `product_variant_stocks` SET current_stock -= Qty, reserved_stock -= Qty
    InventoryRepo->>DB: INSERT INTO `inventory_ledgers` (-Qty, "Order #123")
    InventoryRepo->>DB: COMMIT
    DB-->>Buyer: Order Confirmed
```

## 3. Edge Cases

*   **Order Cancellation/Timeout**:
    *   If the user doesn't pay within the time limit, `ReleaseStock` is called.
    *   `ReservedStock` -= Qty.
*   **Restocking**:
    *   Seller adds more stock. `CurrentStock` += Qty.
