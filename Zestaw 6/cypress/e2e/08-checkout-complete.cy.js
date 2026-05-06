describe('Checkout - finalizacja', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien ukonczyc pelny proces checkout', () => {
    cy.get('[data-test="add-to-cart-sauce-labs-backpack"]').click()
    cy.get('.shopping_cart_link').click()
    cy.get('[data-test="checkout"]').click()
    cy.get('[data-test="firstName"]').type('Jan')
    cy.get('[data-test="lastName"]').type('Kowalski')
    cy.get('[data-test="postalCode"]').type('30-001')
    cy.get('[data-test="continue"]').click()
    cy.url().should('include', '/checkout-step-two.html')
    cy.get('.summary_total_label').should('be.visible')
    cy.get('[data-test="finish"]').click()
    cy.url().should('include', '/checkout-complete.html')
    cy.get('.complete-header').should('contain', 'Thank you for your order')
  })
})
