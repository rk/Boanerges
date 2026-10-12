import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../../wayfinder'
/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
export const show = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(options),
    method: 'get',
})

show.definition = {
    methods: ["get","head"],
    url: '/settings/study',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
show.url = (options?: RouteQueryOptions) => {
    return show.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
show.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: show.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
show.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: show.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
const showForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::show
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: show.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

show.form = showForm

/**
* @see \App\Http\Controllers\SettingsController::update
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
export const update = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: update.url(options),
    method: 'patch',
})

update.definition = {
    methods: ["patch"],
    url: '/settings/study',
} satisfies RouteDefinition<["patch"]>

/**
* @see \App\Http\Controllers\SettingsController::update
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
update.url = (options?: RouteQueryOptions) => {
    return update.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::update
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
update.patch = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: update.url(options),
    method: 'patch',
})

/**
* @see \App\Http\Controllers\SettingsController::update
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
const updateForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\SettingsController::update
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
updateForm.patch = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: update.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

update.form = updateForm

const study = {
    show: Object.assign(show, show),
    update: Object.assign(update, update),
}

export default study