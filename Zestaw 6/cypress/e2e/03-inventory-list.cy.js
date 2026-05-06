describe('Produkty - lista i sortowanie nazw', () => {
  beforeEach(() => {
    cy.login()
  })

  it('powinien wyswietlic liste produktow', () => {
    cy.get('.inventory_list').should('be.visible')
    cy.get('.inventory_item').should('have.length', 6)
    cy.get('.inventory_item').first().find('.inventory_item_name').should('not.be.empty')
    cy.get('.inventory_item').first().find('.inventory_item_price').should('not.be.empty')
    cy.get('.inventory_item').first().find('.inventory_item_img').should('be.visible')
    cy.get('.inventory_item').first().find('button').should('be.visible')
  })

  it('powinien sortowac produkty po nazwie A-Z', () => {
    cy.get('[data-test="product-sort-container"]').select('az')
    cy.get('.inventory_item_name').first().should('contain', 'Sauce Labs Backpack')
  })

  it('powinien sortowac produkty po nazwie Z-A', () => {
    cy.get('[data-test="product-sort-container"]').select('za')
    cy.get('.inventory_item_name').first().invoke('text').then((text) => {
      expect(text.trim()).to.equal('Test.allTheThings() T-Shirt (Red)')
    })
  })
})
