<?php

namespace App\Http\Resources\Bible;

use App\Data\TranslationConfig;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin TranslationConfig */
class TranslationResource extends JsonResource
{
    /**
     * @return array{id: string, name: string, abbrev: string, bundled: bool, about: string|null, install_status: string|null, install_step: string|null, install_error: string|null}
     */
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'abbrev' => $this->abbrev,
            'about' => $this->about ?? "https://crosswire.org/sword/modules/ModInfo.jsp?modName={$this->abbrev}",
            'bundled' => $this->bundled,
            'install_status' => $this->installStatus?->value,
            'install_step' => $this->installStep?->value,
            'install_error' => $this->installError,
        ];
    }
}
