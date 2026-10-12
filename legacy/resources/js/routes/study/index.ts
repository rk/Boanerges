import { queryParams, type RouteQueryOptions, type RouteDefinition, type RouteFormDefinition } from './../../wayfinder'
import printers from './printers'
/**
* @see \App\Http\Controllers\StudyPrintController::print
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
export const print = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: print.url(options),
    method: 'post',
})

print.definition = {
    methods: ["post"],
    url: '/study/print',
} satisfies RouteDefinition<["post"]>

/**
* @see \App\Http\Controllers\StudyPrintController::print
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
print.url = (options?: RouteQueryOptions) => {
    return print.definition.url + queryParams(options)
}

/**
* @see \App\Http\Controllers\StudyPrintController::print
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
print.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
    url: print.url(options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\StudyPrintController::print
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
const printForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: print.url(options),
    method: 'post',
})

/**
* @see \App\Http\Controllers\StudyPrintController::print
* @see app/Http/Controllers/StudyPrintController.php:19
* @route '/study/print'
*/
printForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
    action: print.url(options),
    method: 'post',
})

print.form = printForm

const study = {
    printers: Object.assign(printers, printers),
    print: Object.assign(print, print),
}

export default study