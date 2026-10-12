import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../../../wayfinder'
/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
export const showReadability = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: showReadability.url(options),
    method: 'get',
})

showReadability.definition = {
    methods: ["get","head"],
    url: '/settings/readability',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
showReadability.url = (options?: RouteQueryOptions) => {
    return showReadability.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
showReadability.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: showReadability.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
showReadability.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: showReadability.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
const showReadabilityForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showReadability.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
showReadabilityForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showReadability.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showReadability
* @see app/Http/Controllers/SettingsController.php:14
* @route '/settings/readability'
*/
showReadabilityForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showReadability.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

showReadability.form = showReadabilityForm

/**
* @see \App\Http\Controllers\SettingsController::updateReadability
* @see app/Http/Controllers/SettingsController.php:21
* @route '/settings/readability'
*/
export const updateReadability = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: updateReadability.url(options),
    method: 'patch',
})

updateReadability.definition = {
    methods: ["patch"],
    url: '/settings/readability',
} satisfies RouteDefinition<["patch"]>

/**
* @see \App\Http\Controllers\SettingsController::updateReadability
* @see app/Http/Controllers/SettingsController.php:21
* @route '/settings/readability'
*/
updateReadability.url = (options?: RouteQueryOptions) => {
    return updateReadability.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::updateReadability
* @see app/Http/Controllers/SettingsController.php:21
* @route '/settings/readability'
*/
updateReadability.patch = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: updateReadability.url(options),
    method: 'patch',
})

/**
* @see \App\Http\Controllers\SettingsController::updateReadability
* @see app/Http/Controllers/SettingsController.php:21
* @route '/settings/readability'
*/
const updateReadabilityForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: updateReadability.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\SettingsController::updateReadability
* @see app/Http/Controllers/SettingsController.php:21
* @route '/settings/readability'
*/
updateReadabilityForm.patch = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: updateReadability.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

updateReadability.form = updateReadabilityForm

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
export const showStudy = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: showStudy.url(options),
    method: 'get',
})

showStudy.definition = {
    methods: ["get","head"],
    url: '/settings/study',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showStudy.url = (options?: RouteQueryOptions) => {
    return showStudy.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showStudy.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: showStudy.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showStudy.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: showStudy.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
const showStudyForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showStudy.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showStudyForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showStudy.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\SettingsController::showStudy
* @see app/Http/Controllers/SettingsController.php:30
* @route '/settings/study'
*/
showStudyForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: showStudy.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

showStudy.form = showStudyForm

/**
* @see \App\Http\Controllers\SettingsController::updateStudy
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
export const updateStudy = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: updateStudy.url(options),
    method: 'patch',
})

updateStudy.definition = {
    methods: ["patch"],
    url: '/settings/study',
} satisfies RouteDefinition<["patch"]>

/**
* @see \App\Http\Controllers\SettingsController::updateStudy
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
updateStudy.url = (options?: RouteQueryOptions) => {
    return updateStudy.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\SettingsController::updateStudy
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
updateStudy.patch = (options?: RouteQueryOptions): RouteDefinition<'patch'> => ({
    url: updateStudy.url(options),
    method: 'patch',
})

/**
* @see \App\Http\Controllers\SettingsController::updateStudy
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
const updateStudyForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: updateStudy.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\SettingsController::updateStudy
* @see app/Http/Controllers/SettingsController.php:37
* @route '/settings/study'
*/
updateStudyForm.patch = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: updateStudy.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'PATCH',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

updateStudy.form = updateStudyForm

const SettingsController = { showReadability, updateReadability, showStudy, updateStudy }

export default SettingsController