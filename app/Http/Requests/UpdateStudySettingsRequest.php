<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;
use Illuminate\Validation\Validator;

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

                return;
            }

            $nonBible = [];

            foreach ($columns as $index => $column) {
                if ($column === 'bible-secondary') {
                    continue;
                }

                if (in_array($column, $nonBible, true)) {
                    $validator->errors()->add("columns.{$index}", 'Duplicate column type is not allowed.');

                    return;
                }

                $nonBible[] = $column;
            }

            $bibleSlots = array_keys(array_filter($columns, fn($column) => $column === 'bible-secondary'));

            if (count($bibleSlots) === 2) {
                $translations = [
                    0 => (string) $this->input('translationBId'),
                    1 => (string) $this->input('translationCId'),
                ];

                if ($translations[$bibleSlots[0]] === $translations[$bibleSlots[1]]) {
                    $validator->errors()->add('translationCId', 'Secondary Bible columns must use different translations.');
                }
            }
        });
    }
}
