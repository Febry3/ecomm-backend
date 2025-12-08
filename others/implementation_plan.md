# Implementation Plan

Based on the database schema, here is a week-by-week plan to implement the features.

## Phase 1: Core Foundation & Users

### Week 1: Authentication & User Management
**Focus**: Secure access and user profile management.
- [X] **Auth System**:
    - Implement `AuthProvider` (Login/Register with Email/Password or OAuth).
    - Implement `RefreshToken` flow (Issue, Validate, Revoke).
    - Middleware for Role-based access (User vs Seller vs Admin).
- [X] **User Profile**:
    - CRUD for `User` (Profile update, Avatar).
- [X] **Address Management** (High Priority):
    - CRUD for `Address`.
    - **Task**: Implement "Main Address" logic (Ensure only one `IsDefault` per user).

## Phase 2: Supply Side (Sellers & Products)

### Week 2: Seller Onboarding & Product Catalog
**Focus**: Enabling sellers to list products.
- [ ] **Seller Module**:
    - Seller Registration (`Seller` table).
    - Shop Profile management (Logo, Description).
- [ ] **Category Management**:
    - Tree structure for `Category` (Parent-Child).
- [ ] **Product Management**:
    - CRUD for `Product` (Basic info).
    - CRUD for `ProductVariant` (SKU, Price).
    - `ProductImage` upload and management.
- [ ] **Inventory System**:
    - **Entities**:
        - Update `ProductVariant` to has-one `ProductVariantStock`.
        - Update `ProductVariantStock` to belong-to `ProductVariant`.
    - **Repository**:
        - Create `InventoryRepository` (GetStock, UpdateStock).
        - Implement Transactional updates (Stock + Ledger).
    - **Usecase**:
        - `AddStock` (Restocking).
        - `ReserveStock` (Order Placement).
        - `ReleaseStock` (Order Cancellation).
        - `DeductStock` (Order Fulfillment).

## Phase 3: Discovery & Shopping Core

### Week 3: Discovery & Engagement
**Focus**: Helping users find products and engage.
- [ ] **Search & Filter**:
    - Product listing with filters (Category, Price, etc.).
- [ ] **Social Features**:
    - `UserFavorite` (Wishlist).
    - `ProductReview` (Ratings & Text).
    - `SellerReview` (Store ratings).

### Week 4: Cart & Order Processing
**Focus**: The checkout flow.
- [ ] **Shopping Cart**:
    - `Cart` and `CartItem` management (Redis or DB).
    - Session management for guest vs logged-in users.
- [ ] **Checkout**:
    - `Order` creation.
    - `OrderItem` snapshotting (Price at purchase).
    - `OrderShippingDetail` capture.
- [ ] **Payment Integration**:
    - `Payment` record creation.
    - Mock payment gateway integration (or real one like Midtrans/Stripe).

## Phase 4: Advanced Features (The "Startup" Value)

### Week 5: Group Buy System
**Focus**: The core differentiator feature.
- [ ] **Group Buy Logic**:
    - Create `GroupBuySession` (Organizer, Product, Discount).
    - Join Session (`GroupBuyParticipant`).
    - Logic to check `MinParticipants` and `ExpiresAt`.
- [ ] **Order Integration**:
    - Link `Order` to `GroupBuySession`.
    - Handle `OrderAdjustment` for group buy discounts.

### Week 6: Financials & Seller Operations
**Focus**: Money movement and order fulfillment.
- [ ] **Order Fulfillment**:
    - Seller Dashboard to view and update `Order` status (Shipped, Delivered).
- [ ] **Commissions & Payouts**:
    - Calculate `SellerCommission` per order.
    - Implement `SellerPayout` request and processing flow.
- [ ] **Coupons**:
    - Implement `Coupon` logic and `OrderAdjustment`.

## Phase 7: Polish & Launch

### Week 7: Optimization & Testing
**Focus**: Stability and performance.
- [ ] **Performance**:
    - Database indexing optimization.
    - Caching for Products/Categories.
- [ ] **Testing**:
    - Integration tests for Order Flow and Group Buy logic.
- [ ] **Final Polish**:
    - UI/UX improvements.
    - SEO optimizations.
