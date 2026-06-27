<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;
use Illuminate\Validation\Validator;

class PrintStudyRequest extends FormRequest
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
            'includeUserWork' => ['required', 'boolean'],
            'printerName' => ['nullable', 'string', 'max:255'],
            'columnCount' => ['required', 'integer', Rule::in([1, 2, 3])],
            'columns' => ['present', 'array'],
            'columns.*' => ['string', Rule::in(['bible-secondary', 'notes', 'scribe', 'search', 'cross-references'])],
            'bookId' => ['required', 'string', 'max:10'],
            'chapter' => ['required', 'integer', 'min:1', 'max:150'],
            'translationId' => ['required', 'string', 'max:10'],
            'translationBId' => ['required', 'string', 'max:10'],
            'translationCId' => ['required', 'string', 'max:10'],
        ];
    }

    public function withValidator(Validator $validator): void
    {
        $validator->after(function (Validator $validator): void {
            $columnCount = (int) $this->input('columnCount', 1);
            $columns = $this->input('columns', []);

            if (! is_array($columns)) {
                return;
            }

            if (count($columns) !== max(0, $columnCount - 1)) {
                $validator->errors()->add('columns', 'Column slots must match column count.');
            }
        });
    }

    /**
     * @return array{
     *     columnCount: int,
     *     columns: list<string>,
     *     bookId: string,
     *     chapter: int,
     *     translationId: string,
     *     translationBId: string,
     *     translationCId: string
     * }
     */
    public function studySettings(): array
    {
        return [
            'columnCount' => (int) $this->validated('columnCount'),
            'columns' => array_values($this->validated('columns')),
            'bookId' => (string) $this->validated('bookId'),
            'chapter' => (int) $this->validated('chapter'),
            'translationId' => (string) $this->validated('translationId'),
            'translationBId' => (string) $this->validated('translationBId'),
            'translationCId' => (string) $this->validated('translationCId'),
        ];
    }
}
