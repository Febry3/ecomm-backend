# Group Buy System - "One Night" Speedrun Plan

This plan is designed to implement the core Group Buy functionality **and its dependencies** (Seller & Product) in a single coding session. We will focus on the *minimum viable* implementation to get the flow working.

## Phase 0: Fix Environment (5 mins)
- [ ] **Docker Port Conflict**: Fix the `bind: address already in use` error by changing the host port for Postgres in `docker-compose.yml` (e.g., `5433:5432`).

## Phase 1: Seller & Product (Dependencies) (1 hour)
Since Group Buy requires Products, and Products require Sellers, we must build these first.
- [ ] **Seller Module**:
    - **Repo**: `SellerRepository` (Create, FindByUserID, Update).
    - **Usecase**: 
        - `RegisterSeller` (upgrades User to Seller).
        - `GetSeller` (Get seller profile by ID or UserID).
        - `UpdateSeller` (Update shop name, description, logo).
    - **Handler**: 
        - `POST /sellers` (Register).
        - `GET /sellers/me` (Get My Profile).
        - `PUT /sellers/me` (Update Profile).
- [ ] **Product Module** (Simplified):
    - **Repo**: `ProductRepository` (Create, FindByID).
    - **Usecase**: `ProductUsecase` (CreateProduct - creates Product + 1 Variant).
    - **Handler**: `POST /products` (Create).
    - *Note*: We will skip complex Category trees and multiple images for tonight. Just basic Product + Variant.

## Phase 2: Group Buy Core (1.5 hours)
- [ ] **Preparation**:
    - Verify `group_buy_sessions` and `group_buy_participants` tables.
    - Create DTOs (`CreateGroupBuyRequest`, `JoinGroupBuyRequest`).
- [ ] **Repository Layer**:
    - `GroupBuyRepository`: Create, FindByID, AddParticipant, CountParticipants.
- [ ] **Usecase Layer**:
    - `CreateSession`: Validate params, create session.
    - `JoinSession`: Check expiry, max participants, existing participation.
- [ ] **Delivery Layer**:
    - `POST /group-buys`: Create a session for a product.
    - `POST /group-buys/:id/join`: Join a session.

## Phase 3: Order Integration (45 mins)
- [ ] **Order Creation**:
    - Update `CreateOrder` to accept `GroupBuySessionID`.
    - Validate session status and user participation.
    - Apply `DiscountPercentage` to the price.
    - Record `OrderAdjustment`.

## Phase 4: Testing (Remaining Time)
- [ ] **End-to-End Flow**:
    1. Register as Seller.
    2. Create a Product.
    3. Create a Group Buy Session for that Product.
    4. (As another User) Join the Session.
    5. Place an Order linked to the Session.
    6. Verify Discount.


## 7. Bonus (If time permits)
- [ ] **Expiry Check**: A simple endpoint `POST /group-buys/check-expiry` that iterates active sessions and marks them `cancelled` if `ExpiresAt` < Now and `Participants` < `Min`.
