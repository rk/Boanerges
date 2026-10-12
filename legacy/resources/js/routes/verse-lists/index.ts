import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../wayfinder'
/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

index.definition = {
    methods: ["get","head"],
    url: '/verse-lists',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
index.url = (options?: RouteQueryOptions) => {
    return index.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: index.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
*/
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::index
* @see app/Http/Controllers/VerseListController.php:12
* @route '/verse-lists'
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
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
export const show = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

show.definition = {
    methods: ["get","head"],
    url: '/verse-lists/{id}',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
show.url = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { id: args }
    }

    if (Array.isArray(args)) {
        args = {
            id: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        id: args.id,
    }

    return show.definition.url
            .replace('{id}', parsedArgs.id.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
show.get = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
show.head = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: show.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
const showForm = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
showForm.get = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\VerseListController::show
* @see app/Http/Controllers/VerseListController.php:19
* @route '/verse-lists/{id}'
*/
showForm.head = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
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
* @see \App\Http\Controllers\VerseListController::store
* @see app/Http/Controllers/VerseListController.php:30
* @route '/verse-lists'
*/
export const store = (options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: store.url(options),
    method: 'put',
})

store.definition = {
    methods: ["put"],
    url: '/verse-lists',
} satisfies RouteDefinition<["put"]>

/**
* @see \App\Http\Controllers\VerseListController::store
* @see app/Http/Controllers/VerseListController.php:30
* @route '/verse-lists'
*/
store.url = (options?: RouteQueryOptions) => {
    return store.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\VerseListController::store
* @see app/Http/Controllers/VerseListController.php:30
* @route '/verse-lists'
*/
store.put = (options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: store.url(options),
    method: 'put',
})

/**
* @see \App\Http\Controllers\VerseListController::store
* @see app/Http/Controllers/VerseListController.php:30
* @route '/verse-lists'
*/
const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: store.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\VerseListController::store
* @see app/Http/Controllers/VerseListController.php:30
* @route '/verse-lists'
*/
storeForm.put = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: store.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

store.form = storeForm

/**
* @see \App\Http\Controllers\VerseListController::update
* @see app/Http/Controllers/VerseListController.php:37
* @route '/verse-lists/{id}'
*/
export const update = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: update.url(args, options),
    method: 'put',
})

update.definition = {
    methods: ["put"],
    url: '/verse-lists/{id}',
} satisfies RouteDefinition<["put"]>

/**
* @see \App\Http\Controllers\VerseListController::update
* @see app/Http/Controllers/VerseListController.php:37
* @route '/verse-lists/{id}'
*/
update.url = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { id: args }
    }

    if (Array.isArray(args)) {
        args = {
            id: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        id: args.id,
    }

    return update.definition.url
            .replace('{id}', parsedArgs.id.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\VerseListController::update
* @see app/Http/Controllers/VerseListController.php:37
* @route '/verse-lists/{id}'
*/
update.put = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
    url: update.url(args, options),
    method: 'put',
})

/**
* @see \App\Http\Controllers\VerseListController::update
* @see app/Http/Controllers/VerseListController.php:37
* @route '/verse-lists/{id}'
*/
const updateForm = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\VerseListController::update
* @see app/Http/Controllers/VerseListController.php:37
* @route '/verse-lists/{id}'
*/
updateForm.put = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PUT',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

update.form = updateForm

/**
* @see \App\Http\Controllers\VerseListController::destroy
* @see app/Http/Controllers/VerseListController.php:45
* @route '/verse-lists/{id}'
*/
export const destroy = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: destroy.url(args, options),
    method: 'delete',
})

destroy.definition = {
    methods: ["delete"],
    url: '/verse-lists/{id}',
} satisfies RouteDefinition<["delete"]>

/**
* @see \App\Http\Controllers\VerseListController::destroy
* @see app/Http/Controllers/VerseListController.php:45
* @route '/verse-lists/{id}'
*/
destroy.url = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { id: args }
    }

    if (Array.isArray(args)) {
        args = {
            id: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        id: args.id,
    }

    return destroy.definition.url
            .replace('{id}', parsedArgs.id.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\VerseListController::destroy
* @see app/Http/Controllers/VerseListController.php:45
* @route '/verse-lists/{id}'
*/
destroy.delete = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: destroy.url(args, options),
    method: 'delete',
})

/**
* @see \App\Http\Controllers\VerseListController::destroy
* @see app/Http/Controllers/VerseListController.php:45
* @route '/verse-lists/{id}'
*/
const destroyForm = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: destroy.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\VerseListController::destroy
* @see app/Http/Controllers/VerseListController.php:45
* @route '/verse-lists/{id}'
*/
destroyForm.delete = (args: { id: string | number } | [id: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: destroy.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

destroy.form = destroyForm

const verseLists = {
    index: Object.assign(index, index),
    show: Object.assign(show, show),
    store: Object.assign(store, store),
    update: Object.assign(update, update),
    destroy: Object.assign(destroy, destroy),
}

export default verseLists