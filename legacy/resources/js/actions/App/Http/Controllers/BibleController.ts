import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../../wayfinder'
/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
export const translations = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: translations.url(options),
    method: 'get',
})

translations.definition = {
    methods: ["get","head"],
    url: '/bible/translations',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
translations.url = (options?: RouteQueryOptions) => {
    return translations.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
translations.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: translations.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
translations.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: translations.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
const translationsForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: translations.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
translationsForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: translations.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::translations
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
translationsForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: translations.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

translations.form = translationsForm

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
export const catalog = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: catalog.url(options),
    method: 'get',
})

catalog.definition = {
    methods: ["get","head"],
    url: '/bible/translations/catalog',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.url = (options?: RouteQueryOptions) => {
    return catalog.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: catalog.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
const catalogForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalogForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalogForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

catalog.form = catalogForm

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
export const install = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: install.url(args, options),
    method: 'post',
})

install.definition = {
    methods: ["post"],
    url: '/bible/translations/{module}/install',
} satisfies RouteDefinition<["post"]>

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
install.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return install.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
install.post = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: install.url(args, options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
const installForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: install.url(args, options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
installForm.post = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: install.url(args, options),
    method: 'post',
})

install.form = installForm

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
export const installStatus = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: installStatus.url(args, options),
    method: 'get',
})

installStatus.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{module}/install-status',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return installStatus.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.get = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.head = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: installStatus.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
const installStatusForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatusForm.get = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatusForm.head = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

installStatus.form = installStatusForm

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
export const uninstall = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: uninstall.url(args, options),
    method: 'delete',
})

uninstall.definition = {
    methods: ["delete"],
    url: '/bible/translations/{module}',
} satisfies RouteDefinition<["delete"]>

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstall.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return uninstall.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstall.delete = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: uninstall.url(args, options),
    method: 'delete',
})

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
const uninstallForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: uninstall.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstallForm.delete = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: uninstall.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

uninstall.form = uninstallForm

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
export const search = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: search.url(options),
    method: 'get',
})

search.definition = {
    methods: ["get","head"],
    url: '/bible/search',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
search.url = (options?: RouteQueryOptions) => {
    return search.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
search.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: search.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
search.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: search.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
const searchForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: search.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
searchForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: search.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::search
* @see app/Http/Controllers/BibleController.php:110
* @route '/bible/search'
*/
searchForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: search.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

search.form = searchForm

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
export const crossReferences = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: crossReferences.url(options),
    method: 'get',
})

crossReferences.definition = {
    methods: ["get","head"],
    url: '/bible/cross-references',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
crossReferences.url = (options?: RouteQueryOptions) => {
    return crossReferences.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
crossReferences.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: crossReferences.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
crossReferences.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: crossReferences.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
const crossReferencesForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: crossReferences.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
crossReferencesForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: crossReferences.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::crossReferences
* @see app/Http/Controllers/BibleController.php:127
* @route '/bible/cross-references'
*/
crossReferencesForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: crossReferences.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

crossReferences.form = crossReferencesForm

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
export const books = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: books.url(args, options),
    method: 'get',
})

books.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{translation}/books',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
books.url = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { translation: args }
    }

    if (Array.isArray(args)) {
        args = {
            translation: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        translation: args.translation,
    }

    return books.definition.url
            .replace('{translation}', parsedArgs.translation.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
books.get = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: books.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
books.head = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: books.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
const booksForm = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: books.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
booksForm.get = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: books.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::books
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
booksForm.head = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: books.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

books.form = booksForm

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
export const chapter = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: chapter.url(args, options),
    method: 'get',
})

chapter.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{translation}/books/{book}/chapters/{chapter}',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
chapter.url = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions) => {
    if (Array.isArray(args)) {
        args = {
            translation: args[0],
            book: args[1],
            chapter: args[2],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        translation: args.translation,
        book: args.book,
        chapter: args.chapter,
    }

    return chapter.definition.url
            .replace('{translation}', parsedArgs.translation.toString())
            .replace('{book}', parsedArgs.book.toString())
            .replace('{chapter}', parsedArgs.chapter.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
chapter.get = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: chapter.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
chapter.head = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: chapter.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
const chapterForm = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: chapter.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
chapterForm.get = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: chapter.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::chapter
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
chapterForm.head = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: chapter.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

chapter.form = chapterForm

const BibleController = { translations, catalog, install, installStatus, uninstall, search, crossReferences, books, chapter }

export default BibleController