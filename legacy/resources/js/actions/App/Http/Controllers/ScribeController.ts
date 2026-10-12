import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../../wayfinder'
/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
export const show = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

show.definition = {
    methods: ["get","head"],
    url: '/scribe/{book}/chapters/{chapter}',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
show.url = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions) => {
    if (Array.isArray(args)) {
        args = {
            book: args[0],
            chapter: args[1],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        book: args.book,
        chapter: args.chapter,
    }

    return show.definition.url
            .replace('{book}', parsedArgs.book.toString())
            .replace('{chapter}', parsedArgs.chapter.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
show.get = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
show.head = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: show.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
const showForm = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
showForm.get = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\ScribeController::show
* @see app/Http/Controllers/ScribeController.php:12
* @route '/scribe/{book}/chapters/{chapter}'
*/
showForm.head = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

show.form = showForm

/**
* @see \App\Http\Controllers\ScribeController::update
* @see app/Http/Controllers/ScribeController.php:26
* @route '/scribe/{book}/chapters/{chapter}'
*/
export const update = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: update.url(args, options),
    method: 'put',
})

update.definition = {
    methods: ["put"],
    url: '/scribe/{book}/chapters/{chapter}',
} satisfies RouteDefinition<["put"]>

/**
* @see \App\Http\Controllers\ScribeController::update
* @see app/Http/Controllers/ScribeController.php:26
* @route '/scribe/{book}/chapters/{chapter}'
*/
update.url = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions) => {
    if (Array.isArray(args)) {
        args = {
            book: args[0],
            chapter: args[1],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        book: args.book,
        chapter: args.chapter,
    }

    return update.definition.url
            .replace('{book}', parsedArgs.book.toString())
            .replace('{chapter}', parsedArgs.chapter.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\ScribeController::update
* @see app/Http/Controllers/ScribeController.php:26
* @route '/scribe/{book}/chapters/{chapter}'
*/
update.put = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: update.url(args, options),
    method: 'put',
})

/**
* @see \App\Http\Controllers\ScribeController::update
* @see app/Http/Controllers/ScribeController.php:26
* @route '/scribe/{book}/chapters/{chapter}'
*/
const updateForm = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\ScribeController::update
* @see app/Http/Controllers/ScribeController.php:26
* @route '/scribe/{book}/chapters/{chapter}'
*/
updateForm.put = (args: { book: string | number, chapter: string | number } | [book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

update.form = updateForm

const ScribeController = { show, update }

export default ScribeController