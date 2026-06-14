<?php

namespace App\Http\Resources\Bible;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class ChapterResource extends JsonResource
{
    /**
     * @return array{
     *     book: string,
     *     bookAbbrev: string,
     *     chapter: int,
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var array{
         *     book: string,
         *     bookAbbrev: string,
         *     chapter: int,
         *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
         * } $chapter
         */
        $chapter = $this->resource;

        return $chapter;
    }
}
