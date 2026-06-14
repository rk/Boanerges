<?php

namespace App\Http\Resources\Bible;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class BookResource extends JsonResource
{
    /**
     * @return array{
     *     id: string,
     *     name: string,
     *     abbrev: string,
     *     testament: string,
     *     chapters: int,
     *     firstChapter: int,
     *     lastChapter: int
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var array{
         *     id: string,
         *     name: string,
         *     abbrev: string,
         *     testament: string,
         *     chapters: int,
         *     firstChapter: int,
         *     lastChapter: int
         * } $book
         */
        $book = $this->resource;

        return $book;
    }
}
