<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;

class StoreVerseListRequest extends FormRequest
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
            'title' => ['required', 'string', 'max:255'],
            'entries' => ['present', 'array'],
            'entries.*.bookId' => ['required', 'string', 'max:10'],
            'entries.*.chapter' => ['required', 'integer', 'min:1', 'max:150'],
            'entries.*.verse' => ['required', 'integer', 'min:1'],
        ];
    }
}
