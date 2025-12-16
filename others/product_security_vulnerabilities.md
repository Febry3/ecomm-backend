# Product Service Security Vulnerabilities

This document outlines security vulnerabilities identified in the Product Service module.

---

## 1. IDOR (Insecure Direct Object Reference) in Variant Updates

**Severity:** 🔴 HIGH

**Location:** `internal/usecase/product_usecase.go` - `UpdateProduct()` (lines 169-181)

**Issue:** When updating product variants, the code does not verify that the variant belongs to the product or seller.

```go
// Current vulnerable code
if err := p.variantRepo.UpdateProductVariant(txCtx, variant, v.ID); err != nil {
    // No check if v.ID belongs to the product being updated
}
```

**Attack Scenario:**
1. Attacker is a seller with their own products
2. Attacker finds variant ID of another seller's product
3. Attacker sends update request with their product ID but another seller's variant ID
4. Variant is updated without ownership verification

**Fix:** Verify variant ownership before update:
```go
// Check variant belongs to this product
existingVariant, err := p.variantRepo.GetProductVariant(txCtx, v.ID)
if err != nil || existingVariant.ProductID != productID {
    return errors.New("variant not found or unauthorized")
}
```

---

## 2. Missing Seller Role Validation in Some Endpoints

**Severity:** 🟡 MEDIUM

**Location:** `internal/delivery/http/product_delivery.go`

**Issue:** `CreateProduct()` validates seller role (lines 36-44), but `GetAllProductsForSeller()`, `GetProductForSeller()`, and `UpdateProduct()` do not.

**Attack Scenario:**
A non-seller user with a manipulated JWT could potentially access seller endpoints.

**Fix:** Add consistent role validation:
```go
if jwt.Role != "seller" || jwt.SellerID == 0 {
    c.JSON(http.StatusForbidden, gin.H{"message": "seller access required"})
    return
}
```

---

## 3. Stock Manipulation Without Ownership Check

**Severity:** 🔴 HIGH

**Location:** `internal/usecase/product_usecase.go` - `UpdateProduct()` (lines 183-201)

**Issue:** Stock updates are performed using variant ID from request without verifying ownership.

```go
stock := &entity.ProductVariantStock{
    ProductVariantID:  v.ID,  // From user request, unvalidated
    // ...
}
if err := p.stockRepo.UpdateStock(txCtx, stock, v.ID); err != nil {
    // Can manipulate any variant's stock
}
```

**Attack Scenario:**
Attacker could set competitor's stock to 0, making products appear "out of stock."

**Fix:** Same as #1 - verify variant ownership before any operations.

---

## 4. Potential Nil Pointer Dereference

**Severity:** 🟡 MEDIUM

**Location:** `internal/usecase/product_usecase.go` - `UpdateProduct()` (line 159)

**Issue:** `product.IsActive` is dereferenced without nil check.

```go
IsActive: *product.IsActive,  // Crash if IsActive is nil
```

**Fix:** Add nil check:
```go
if product.IsActive != nil {
    productEntity.IsActive = *product.IsActive
}
```

---

## 5. DeleteProductVariant Without Authorization

**Severity:** 🔴 HIGH

**Location:** `internal/usecase/product_usecase.go` - `DeleteProductVariant()` (lines 143-145)

**Issue:** Function exists in usecase but has no seller ownership verification:

```go
func (p *ProductUsecase) DeleteProductVariant(ctx context.Context, productVariantID string) error {
    return p.variantRepo.DeleteProductVariant(ctx, productVariantID)
    // No ownership check - any authenticated user could delete any variant
}
```

**Fix:** Add ownership verification:
```go
func (p *ProductUsecase) DeleteProductVariant(ctx context.Context, productVariantID string, sellerID int64) error {
    variant, err := p.variantRepo.GetProductVariant(ctx, productVariantID)
    if err != nil {
        return err
    }
    
    product, err := p.productRepo.GetProductForSeller(ctx, variant.ProductID, sellerID)
    if err != nil {
        return errors.New("unauthorized: variant does not belong to seller")
    }
    
    return p.variantRepo.DeleteProductVariant(ctx, productVariantID)
}
```

---

## 6. Missing UUID Validation

**Severity:** 🟢 LOW

**Location:** Multiple repository files

**Issue:** UUID parameters from user input are passed directly to GORM without validation.

**Attack Scenario:** Malformed UUID strings could cause database errors or unexpected behavior.

**Fix:** Validate UUID format before queries:
```go
import "github.com/google/uuid"

if _, err := uuid.Parse(productID); err != nil {
    return nil, errors.New("invalid product ID format")
}
```

---

## 7. No Rate Limiting

**Severity:** 🟡 MEDIUM

**Location:** All product endpoints

**Issue:** No rate limiting on product CRUD operations could allow:
- Brute-force enumeration of product/variant IDs
- DoS attacks via rapid creation of products

**Fix:** Implement rate limiting middleware.

---

## Summary Table

| # | Vulnerability | Severity | Status |
|---|--------------|----------|--------|
| 1 | IDOR in Variant Updates | HIGH | Open |
| 2 | Missing Role Validation | MEDIUM | Partial |
| 3 | Stock Manipulation | HIGH | Open |
| 4 | Nil Pointer Dereference | MEDIUM | Open |
| 5 | Delete Without Auth | HIGH | Open |
| 6 | Missing UUID Validation | LOW | Open |
| 7 | No Rate Limiting | MEDIUM | Open |

---

## Priority Recommendations

1. **Immediate:** Fix IDOR vulnerabilities (#1, #3, #5)
2. **Short-term:** Add consistent role validation (#2) and nil checks (#4)
3. **Medium-term:** Add UUID validation (#6) and rate limiting (#7)
