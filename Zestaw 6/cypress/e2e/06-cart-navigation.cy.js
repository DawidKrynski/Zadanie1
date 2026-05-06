describe('Koszyk - widok i nawigacja', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien wyswietlic produkty w widoku koszyka', () => {
    cy.get('[data-test="add-to-cart-sauce-labs-backpack"]').click()
    cy.get('.shopping_cart_link').click()
    cy.url().should('include', '/cart.html')
    cy.get('.cart_item').should('have.length', 1)
    cy.get('.inventory_item_name').should('contain', 'Sauce Labs Backpack')
    cy.get('.cart_quantity').should('have.text', '1')
  })

  it('powinien kontynuowac zakupy z koszyka', () => {
    cy.get('.shopping_cart_link').click()
    cy.get('[data-test="continue-shopping"]').click()
    cy.url().should('include', '/inventory.html')
  })
})
