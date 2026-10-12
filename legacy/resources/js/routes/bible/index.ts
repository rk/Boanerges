import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../wayfinder'
import translations from './translations'
import dictionary from './dictionary'
import books from './books'
import chapters from './chapters'
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

const bible = {
    translations: Object.assign(translations, translations),
    search: Object.assign(search, search),
    crossReferences: Object.assign(crossReferences, crossReferences),
    dictionary: Object.assign(dictionary, dictionary),
    books: Object.assign(books, books),
    chapters: Object.assign(chapters, chapters),
}

export default bible