import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../../wayfinder'
/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
export const suggest = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: suggest.url(options),
    method: 'get',
})

suggest.definition = {
    methods: ["get","head"],
    url: '/bible/dictionary/suggest',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
suggest.url = (options?: RouteQueryOptions) => {
    return suggest.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
suggest.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: suggest.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
suggest.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: suggest.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
const suggestForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: suggest.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
suggestForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: suggest.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::suggest
* @see app/Http/Controllers/DictionaryController.php:12
* @route '/bible/dictionary/suggest'
*/
suggestForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: suggest.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

suggest.form = suggestForm

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
export const show = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

show.definition = {
    methods: ["get","head"],
    url: '/bible/dictionary/{word}',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
show.url = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { word: args }
    }

    if (Array.isArray(args)) {
        args = {
            word: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        word: args.word,
    }

    return show.definition.url
            .replace('{word}', parsedArgs.word.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
show.get = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
show.head = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: show.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
const showForm = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
showForm.get = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\DictionaryController::show
* @see app/Http/Controllers/DictionaryController.php:29
* @route '/bible/dictionary/{word}'
*/
showForm.head = (args: { word: string | number } | [word: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

show.form = showForm

const DictionaryController = { suggest, show }

export default DictionaryController