<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    public function up(): void
    {
        if (Schema::getConnection()->getDriverName() !== 'sqlite') {
            Schema::table('translations', function (Blueprint $table): void {
                $table->string('format')->nullable()->change();
            });

            return;
        }

        if (! Schema::hasTable('translations')) {
            return;
        }

        if (Schema::hasTable('translations_format_nullable_old')) {
            $this->recoverFromFailedMigration();

            return;
        }

        if ($this->formatColumnIsNullable()) {
            return;
        }

        DB::statement('PRAGMA foreign_keys = OFF');

        Schema::create('translations_format_nullable_new', function (Blueprint $table): void {
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

        DB::statement(<<<'SQL'
            INSERT INTO translations_format_nullable_new (
                id, abbrev, name, format, versification, about,
                version_string, version_date, copyright, copyright_contact, source,
                install_status, install_step, install_error, bundled, created_at, updated_at
            )
            SELECT
                id, abbrev, name, format, versification, about,
                version_string, version_date, copyright, copyright_contact, source,
                install_status, install_step, install_error, bundled, created_at, updated_at
            FROM translations
        SQL);

        Schema::drop('translations');
        Schema::rename('translations_format_nullable_new', 'translations');

        DB::statement('PRAGMA foreign_keys = ON');
    }

    public function down(): void
    {
        if (Schema::getConnection()->getDriverName() !== 'sqlite') {
            Schema::table('translations', function (Blueprint $table): void {
                $table->string('format')->default('sword')->nullable(false)->change();
            });

            return;
        }

        DB::table('translations')->whereNull('format')->update(['format' => 'osis']);
    }

    private function recoverFromFailedMigration(): void
    {
        if (DB::table('translations')->count() === 0) {
            DB::statement(<<<'SQL'
                INSERT INTO translations (
                    id, abbrev, name, format, versification, about,
                    version_string, version_date, copyright, copyright_contact, source,
                    install_status, install_step, install_error, bundled, created_at, updated_at
                )
                SELECT
                    id, abbrev, name, format, versification, about,
                    version_string, version_date, copyright, copyright_contact, source,
                    install_status, install_step, install_error, bundled, created_at, updated_at
                FROM translations_format_nullable_old
            SQL);
        }

        Schema::drop('translations_format_nullable_old');
    }

    private function formatColumnIsNullable(): bool
    {
        $columns = DB::select('PRAGMA table_info(translations)');

        foreach ($columns as $column) {
            if ($column->name === 'format') {
                return (int) $column->notnull === 0;
            }
        }

        return false;
    }
};
