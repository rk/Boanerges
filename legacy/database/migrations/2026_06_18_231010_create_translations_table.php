<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    public function up(): void
    {
        Schema::create('translations', function (Blueprint $table): void {
            $table->id();
            $table->string('abbrev')->unique();
            $table->string('name');
            $table->string('format')->nullable();
            $table->string('versification')->nullable();
            $table->text('about')->nullable();
            $table->string('version_string')->nullable();
            $table->string('version_date')->nullable();
            $table->text('copyright')->nullable();
            $table->string('copyright_contact')->nullable();
            $table->text('source')->nullable();
            $table->string('install_status')->default('pending');
            $table->string('install_step')->nullable();
            $table->text('install_error')->nullable();
            $table->boolean('bundled')->default(false);
            $table->timestamps();
        });

        Schema::create('import_meta', function (Blueprint $table): void {
            $table->string('key')->primary();
            $table->timestamp('completed_at')->nullable();
        });

        Schema::create('cross_reference_verses', function (Blueprint $table): void {
            $table->unsignedInteger('id')->primary();
            $table->string('book_id');
            $table->unsignedSmallInteger('chapter');
            $table->unsignedSmallInteger('verse');
        });

        Schema::create('cross_references', function (Blueprint $table): void {
            $table->id();
            $table->unsignedInteger('source_verse_id');
            $table->unsignedInteger('rank');
            $table->unsignedInteger('target_start_id');
            $table->unsignedInteger('target_end_id')->nullable();
            $table->index(['source_verse_id', 'rank']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('cross_references');
        Schema::dropIfExists('cross_reference_verses');
        Schema::dropIfExists('import_meta');
        Schema::dropIfExists('translations');
    }
};
