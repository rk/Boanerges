<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    public function up(): void
    {
        Schema::create('dictionary_entries', function (Blueprint $table): void {
            $table->id();
            $table->string('word');
            $table->string('word_key')->index();
            $table->string('part_of_speech', 50)->nullable();
            $table->text('definitions');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('dictionary_entries');
    }
};
