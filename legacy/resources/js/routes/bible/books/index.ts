import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../wayfinder'
/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
export const index = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(args, options),
    method: 'get',
})

index.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{translation}/books',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
index.url = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions) => {
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

    return index.definition.url
            .replace('{translation}', parsedArgs.translation.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
index.get = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
index.head = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: index.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
const indexForm = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
indexForm.get = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:84
* @route '/bible/translations/{translation}/books'
*/
indexForm.head = (args: { translation: string | number } | [translation: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

index.form = indexForm

const books = {
    index: Object.assign(index, index),
}

export default books