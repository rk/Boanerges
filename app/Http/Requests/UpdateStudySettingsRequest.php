<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class UpdateStudySettingsRequest extends FormRequest
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
            'activeView' => ['required', 'string', Rule::in(['bible', 'comparison', 'scribe'])],
            'bookId' => ['required', 'string', 'max:10'],
            'chapter' => ['required', 'integer', 'min:1', 'max:150'],
            'translationId' => ['required', 'string', 'max:10'],
            'translationBId' => ['required', 'string', 'max:10'],
        ];
    }
}
