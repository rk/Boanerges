import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../../../wayfinder'
/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

index.definition = {
    methods: ["get","head"],
    url: '/study/printers',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
index.url = (options?: RouteQueryOptions) => {
    return index.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: index.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\StudyPrintController::index
* @see app/Http/Controllers/StudyPrintController.php:12
* @route '/study/printers'
*/
indexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

index.form = indexForm

/**
* @see \App\Http\Controllers\StudyPrintController::store
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: store.url(options),
    method: 'post',
})

store.definition = {
    methods: ["post"],
    url: '/study/print',
} satisfies RouteDefinition<["post"]>

/**
* @see \App\Http\Controllers\StudyPrintController::store
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
store.url = (options?: RouteQueryOptions) => {
    return store.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\StudyPrintController::store
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: store.url(options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\StudyPrintController::store
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: store.url(options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\StudyPrintController::store
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: store.url(options),
    method: 'post',
})

store.form = storeForm

const StudyPrintController = { index, store }

export default StudyPrintController