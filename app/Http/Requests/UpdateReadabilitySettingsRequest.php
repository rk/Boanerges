<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class UpdateReadabilitySettingsRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'fontSize' => ['required', 'integer', 'min:14', 'max:24'],
            'lineHeight' => ['required', 'numeric', 'min:1.4', 'max:2'],
            'theme' => ['required', 'string', Rule::in(['light', 'dark'])],
            'fontFamily' => ['required', 'string', Rule::in(['sans-serif', 'serif'])],
        ];
    }
}
