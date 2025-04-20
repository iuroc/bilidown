"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
import('gede-book-api').then(async ({ Book }) => {
    console.log((await Book.getCategories()).map(i => i.name));
});
