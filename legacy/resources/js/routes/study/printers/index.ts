import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../../wayfinder'
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

const printers = {
    index: Object.assign(index, index),
}

export default printers