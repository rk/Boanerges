import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../wayfinder'
/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
export const show = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

show.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{translation}/books/{book}/chapters/{chapter}',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
show.url = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions) => {
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

    return show.definition.url
            .replace('{translation}', parsedArgs.translation.toString())
            .replace('{book}', parsedArgs.book.toString())
            .replace('{chapter}', parsedArgs.chapter.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
show.get = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
show.head = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: show.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
const showForm = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
showForm.get = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::show
* @see app/Http/Controllers/BibleController.php:96
* @route '/bible/translations/{translation}/books/{book}/chapters/{chapter}'
*/
showForm.head = (args: { translation: string | number, book: string | number, chapter: string | number } | [translation: string | number, book: string | number, chapter: string | number ], options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

show.form = showForm

const chapters = {
    show: Object.assign(show, show),
}

export default chapters