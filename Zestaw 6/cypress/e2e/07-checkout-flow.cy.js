describe('Checkout', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien przejsc do formularza checkout', () => {
    cy.get('[data-test="add-to-cart-sauce-labs-backpack"]').click()
    cy.get('.shopping_cart_link').click()
    cy.get('[data-test="checkout"]').click()
    cy.url().should('include', '/checkout-step-one.html')
    cy.get('[data-test="firstName"]').should('be.visible')
    cy.get('[data-test="lastName"]').should('be.visible')
    cy.get('[data-test="postalCode"]').should('be.visible')
  })

  it('powinien wyswietlic blad przy pustym formularzu checkout', () => {
    cy.get('[data-test="add-to-cart-sauce-labs-backpack"]').click()
    cy.get('.shopping_cart_link').click()
    cy.get('[data-test="checkout"]').click()
    cy.get('[data-test="continue"]').click()
    cy.get('[data-test="error"]').should('be.visible')
    cy.get('[data-test="error"]').should('contain', 'First Name is required')
  })
})
