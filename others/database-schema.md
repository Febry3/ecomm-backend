erDiagram
    User ||--o{ AuthProvider : "has"
    User ||--o{ RefreshToken : "has"
    User ||--o{ Address : "has"
    User ||--o{ Cart : "owns"
    User ||--o{ Order : "places"
    User ||--o{ ProductReview : "writes"
    User ||--o{ UserFavorite : "favorites"
    User ||--o{ GroupBuyParticipant : "participates"
    User ||--o| Seller : "becomes"

    Seller ||--o{ Product : "sells"
    Seller ||--o{ GroupBuySession : "organizes"
    Seller ||--o{ SellerReview : "receives"

    Category ||--o{ Category : "parent-child"
    Category ||--o{ Product : "contains"

    Product ||--o{ ProductVariant : "has"
    Product ||--o{ ProductImage : "has"
    Product ||--o{ ProductReview : "receives"

    ProductVariant ||--o{ CartItem : "in"
    ProductVariant ||--o{ OrderItem : "purchased"
    ProductVariant ||--o{ UserFavorite : "favorited"
    ProductVariant ||--o{ InventoryLedger : "tracks"
    ProductVariant ||--o{ GroupBuySession : "featured"
    ProductVariant ||--|| ProductVariantStock : "has"

    Cart ||--o{ CartItem : "contains"

    Order ||--o{ OrderItem : "contains"
    Order ||--|| OrderShippingDetail : "has"
    Order ||--o{ Payment : "has"
    Order ||--o{ OrderAdjustment : "has"
    Order }o--o| GroupBuySession : "linked"
    Order ||--o{ InventoryLedger : "affects"
    Order ||--o{ SellerCommission : "generates"

    GroupBuySession ||--o{ GroupBuyParticipant : "has"

    Seller ||--o{ SellerCommission : "earns"
    Seller ||--o{ SellerPayout : "receives"
    SellerPayout ||--o{ SellerCommission : "includes"

    User {
        int64 ID PK
        string Username UK
        string FirstName
        string LastName
        string PhoneNumber UK
        string Email UK
        string Role "user,seller,admin"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Seller {
        int64 ID PK
        int64 UserID FK,UK
        string StoreName
        string StoreSlug UK
        string Description
        string LogoURL
        string BusinessEmail
        string BusinessPhone
        string Status "pending,approved,suspended"
        bool IsVerified
        decimal AverageRating
        int TotalSales
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    AuthProvider {
        int64 AuthProviderID PK
        int64 UserID FK
        string Provider
        string ProviderId UK
        string Password
    }

    RefreshToken {
        string TokenId PK
        int64 UserID FK
        string TokenHash UK
        bool IsRevoked
        string DeviceInfo
        timestamp ExpiresAt
        timestamp CreatedAt
    }

    Address {
        uuid ID PK
        int64 UserID FK
        string AddressLine1
        string AddressLine2
        string City
        string State
        string PostalCode
        string Country
        bool IsDefault
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Category {
        uuid ID PK
        string Name
        string Slug UK
        string Description
        uuid ParentID FK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Product {
        uuid ID PK
        int64 SellerID FK
        string Title
        string Slug UK
        string Description
        uuid CategoryID FK
        string Badge
        bool IsActive
        string Status "pending,approved,rejected"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    ProductVariant {
        uuid ID PK
        uuid ProductID FK
        string Sku UK
        string Name
        decimal Price
        bool IsActive
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    ProductVariantStock {
        uuid ProductVariantID PK,FK
        int CurrentStock
        int ReservedStock
        int LowStockThreshold
        timestamp LastUpdated
    }

    ProductImage {
        uuid ID PK
        uuid ProductID FK
        string ImageURL
        string AltText
        int DisplayOrder
        bool IsPrimary
        timestamp CreatedAt
    }

    ProductReview {
        uuid ID PK
        uuid ProductID FK
        int64 UserID FK
        int Rating
        text ReviewText
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    SellerReview {
        uuid ID PK
        int64 SellerID FK
        int64 UserID FK
        uuid OrderID FK
        int Rating
        text ReviewText
        int CommunicationRating
        int ShippingRating
        int ProductRating
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    UserFavorite {
        uuid ID PK
        int64 UserID FK
        uuid ProductVariantID FK
        timestamp CreatedAt
    }

    InventoryLedger {
        int64 ID PK
        uuid ProductVariantID FK
        int QuantityChange
        string Reason
        uuid OrderID FK
        timestamp CreatedAt
    }

    Cart {
        uuid ID PK
        int64 UserID FK
        string SessionID UK
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    CartItem {
        uuid ID PK
        uuid CartID FK
        uuid ProductVariantID FK
        int Quantity
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Order {
        uuid ID PK
        string OrderNumber UK
        int64 UserID FK
        uuid GroupBuySessionID FK
        int64 SellerID FK
        decimal Subtotal
        decimal DeliveryCharge
        decimal TotalAmount
        string Status "pending,processing,shipped,delivered,cancelled"
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    OrderItem {
        uuid ID PK
        uuid OrderID FK
        uuid ProductVariantID FK
        int Quantity
        decimal PriceAtPurchase
        decimal TotalPrice
        timestamp CreatedAt
    }

    OrderShippingDetail {
        uuid ID PK
        uuid OrderID FK
        string FullName
        string Phone
        string AddressLine1
        string AddressLine2
        string City
        string State
        string PostalCode
        string Country
    }

    Payment {
        uuid ID PK
        uuid OrderID FK
        decimal Amount
        string Status "pending,succeeded,failed,refunded"
        string PaymentMethod
        string GatewayTransactionID
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    OrderAdjustment {
        uuid ID PK
        uuid OrderID FK
        string Description
        decimal Amount
        string SourceType "coupon,group_buy"
        uuid SourceID FK
        timestamp CreatedAt
    }

    Coupon {
        uuid ID PK
        string Code UK
        string DiscountType "percentage,fixed"
        decimal DiscountValue
        decimal MinPurchaseAmount
        decimal MaxDiscountAmount
        timestamp ValidFrom
        timestamp ValidUntil
        int UsageLimit
        int UsageCount
        bool IsActive
        timestamp CreatedAt
    }

    GroupBuySession {
        uuid ID PK
        string SessionCode UK
        uuid ProductVariantID FK
        int64 OrganizerID FK
        int MinParticipants
        int MaxParticipants
        decimal DiscountPercentage
        string Status "active,completed,cancelled"
        timestamp ExpiresAt
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    GroupBuyParticipant {
        uuid ID PK
        uuid SessionID FK
        int64 UserID FK
        int Quantity
        timestamp JoinedAt
    }

    SellerCommission {
        uuid ID PK
        int64 SellerID FK
        uuid OrderID FK
        uuid OrderItemID FK
        decimal SaleAmount
        decimal CommissionRate
        decimal CommissionAmount
        decimal SellerEarnings
        string Status "pending,paid,held"
        uuid PayoutID FK
        timestamp PaidAt
        timestamp CreatedAt
    }

    SellerPayout {
        uuid ID PK
        int64 SellerID FK
        decimal Amount
        string Status "pending,processing,completed,failed"
        string PaymentMethod "bank_transfer,paypal"
        string TransactionID
        timestamp ProcessedAt
        timestamp CreatedAt
    }