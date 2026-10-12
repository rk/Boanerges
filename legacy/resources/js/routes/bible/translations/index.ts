import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition, applyUrlDefaults } from './../../../wayfinder'
/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

index.definition = {
    methods: ["get","head"],
    url: '/bible/translations',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
index.url = (options?: RouteQueryOptions) => {
    return index.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: index.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
*/
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: index.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::index
* @see app/Http/Controllers/BibleController.php:22
* @route '/bible/translations'
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
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
export const catalog = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: catalog.url(options),
    method: 'get',
})

catalog.definition = {
    methods: ["get","head"],
    url: '/bible/translations/catalog',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.url = (options?: RouteQueryOptions) => {
    return catalog.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalog.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: catalog.url(options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
const catalogForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalogForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url(options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::catalog
* @see app/Http/Controllers/BibleController.php:29
* @route '/bible/translations/catalog'
*/
catalogForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: catalog.url({
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

catalog.form = catalogForm

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
export const install = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: install.url(args, options),
    method: 'post',
})

install.definition = {
    methods: ["post"],
    url: '/bible/translations/{module}/install',
} satisfies RouteDefinition<["post"]>

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
install.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return install.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
install.post = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: install.url(args, options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
const installForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: install.url(args, options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::install
* @see app/Http/Controllers/BibleController.php:40
* @route '/bible/translations/{module}/install'
*/
installForm.post = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: install.url(args, options),
    method: 'post',
})

install.form = installForm

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
export const installStatus = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: installStatus.url(args, options),
    method: 'get',
})

installStatus.definition = {
    methods: ["get","head"],
    url: '/bible/translations/{module}/install-status',
} satisfies RouteDefinition<["get","head"]>

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return installStatus.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.get = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
    url: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatus.head = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
    url: installStatus.url(args, options),
    method: 'head',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
const installStatusForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatusForm.get = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, options),
    method: 'get',
})

/**
* @see \App\Http\Controllers\BibleController::installStatus
* @see app/Http/Controllers/BibleController.php:52
* @route '/bible/translations/{module}/install-status'
*/
installStatusForm.head = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
    action: installStatus.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'HEAD',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'get',
})

installStatus.form = installStatusForm

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
export const uninstall = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: uninstall.url(args, options),
    method: 'delete',
})

uninstall.definition = {
    methods: ["delete"],
    url: '/bible/translations/{module}',
} satisfies RouteDefinition<["delete"]>

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstall.url = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions) => {
    if (typeof args === 'string' || typeof args === 'number') {
        args = { module: args }
    }

    if (Array.isArray(args)) {
        args = {
            module: args[0],
        }
    }

    args = applyUrlDefaults(args)

    const parsedArgs = {
        module: args.module,
    }

    return uninstall.definition.url
            .replace('{module}', parsedArgs.module.toString())
            .replace(/\/+$/, '') + queryParams(options)
}

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstall.delete = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
    url: uninstall.url(args, options),
    method: 'delete',
})

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
const uninstallForm = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: uninstall.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

/**
* @see \App\Http\Controllers\BibleController::uninstall
* @see app/Http/Controllers/BibleController.php:73
* @route '/bible/translations/{module}'
*/
uninstallForm.delete = (args: { module: string | number } | [module: string | number ] | string | number, options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: uninstall.url(args, {
        [options?.mergeQuery ? 'mergeQuery' : 'query']: {
            _method: 'DELETE',
            ...(options?.query ?? options?.mergeQuery ?? {}),
        }
    }),
    method: 'post',
})

uninstall.form = uninstallForm

const translations = {
    index: Object.assign(index, index),
    catalog: Object.assign(catalog, catalog),
    install: Object.assign(install, install),
    installStatus: Object.assign(installStatus, installStatus),
    uninstall: Object.assign(uninstall, uninstall),
}

export default translations