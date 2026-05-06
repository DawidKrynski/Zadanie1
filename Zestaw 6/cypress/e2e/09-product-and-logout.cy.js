describe('Produkt i wylogowanie', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien otworzyc strone szczegolow produktu', () => {
    cy.get('.inventory_item_name').first().click()
    cy.url().should('include', '/inventory-item.html')
    cy.get('.inventory_details_name').should('be.visible')
    cy.get('.inventory_details_price').should('be.visible')
    cy.get('.inventory_details_desc').should('not.be.empty')
  })

  it('powinien wylogowac uzytkownika', () => {
    cy.get('#react-burger-menu-btn').click()
    cy.get('#logout_sidebar_link').should('be.visible')
    cy.get('#logout_sidebar_link').click()
    cy.url().should('eq', 'https://www.saucedemo.com/')
    cy.get('#login-button').should('be.visible')
  })
})
