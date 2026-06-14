<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;

class UpdateScribeChapterRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    /**
     * @return array<string, ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'verses' => ['present', 'array'],
            'verses.*.verse' => ['required', 'integer', 'min:1', 'max:200'],
            'verses.*.text' => ['present', 'string'],
            'verses.*.paragraphStart' => ['sometimes', 'boolean'],
        ];
    }
}
