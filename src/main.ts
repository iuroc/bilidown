import('gede-book-api').then(async ({ Book }) => {
    console.log((await Book.getCategories()).map(i => i.name))
})