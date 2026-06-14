<?php

namespace App\Http\Resources\Bible;

use App\Data\TranslationConfig;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin TranslationConfig */
class TranslationResource extends JsonResource
{
    /**
     * @return array{id: string, name: string, abbrev: string}
     */
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'abbrev' => $this->abbrev,
        ];
    }
}
