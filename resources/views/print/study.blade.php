<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Study Print</title>
    <style>
        @page {
            size: A4 {{ $landscape ? 'landscape' : 'portrait' }};
            margin: 12mm;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: {{ $fontFamily }};
            font-size: {{ $fontSize }}px;
            line-height: {{ $lineHeight }};
            text-align: {{ $justifyText ? 'justify' : 'left' }};
            color: #111;
        }

        .columns {
            display: flex;
            width: 100%;
            min-height: calc(100vh - 24mm);
        }

        .column {
            flex: 1 1 0;
            min-width: 0;
            padding: 0 4mm;
            border-right: 1px solid #ccc;
        }

        .column:last-child {
            border-right: none;
        }

        .column-title {
            font-weight: 600;
            font-size: 0.95em;
            margin-bottom: 0.75em;
            padding-bottom: 0.35em;
            border-bottom: 1px solid #ddd;
        }

        .column-subtitle {
            font-size: 0.75em;
            color: #666;
            margin-bottom: 0.75em;
        }

        .reader p {
            margin-bottom: 0.75em;
        }

        .reader sup {
            font-size: 0.7em;
            color: #555;
            margin-right: 0.15em;
        }

        .notes-content {
            white-space: pre-wrap;
        }

        .lined-block {
            min-height: 70vh;
            background-image: repeating-linear-gradient(
                transparent,
                transparent calc({{ $lineHeight }}em - 1px),
                #ccc calc({{ $lineHeight }}em - 1px),
                #ccc {{ $lineHeight }}em
            );
            background-attachment: local;
        }

        .scribe-verse {
            display: flex;
            align-items: baseline;
            gap: 0.25em;
            min-height: {{ $lineHeight }}em;
            border-bottom: 1px solid #ccc;
            margin-bottom: 0.15em;
        }

        .scribe-verse.paragraph-start {
            margin-top: 0.85em;
        }

        .scribe-verse sup {
            flex-shrink: 0;
            font-size: 0.7em;
            color: #555;
            min-width: 1.25em;
        }

        .scribe-verse .line {
            flex: 1;
            min-height: {{ $lineHeight }}em;
        }

        .message {
            color: #666;
            font-size: 0.9em;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="columns">
        @foreach ($columns as $column)
            <div class="column column-{{ $column['kind'] }}">
                <div class="column-title">{{ $column['label'] }}</div>

                @if (($column['kind'] ?? '') === 'bible')
                    <div class="reader">
                        @php
                            $paragraph = [];
                        @endphp
                        @foreach ($column['verses'] as $verse)
                            @if (($verse['paragraphStart'] ?? false) && $paragraph !== [])
                                <p>
                                    @foreach ($paragraph as $index => $part)
                                        @if ($index > 0)&nbsp;@endif
                                        <sup>{{ $part['number'] }}</sup>{{ $part['text'] }}
                                    @endforeach
                                </p>
                                @php $paragraph = []; @endphp
                            @endif
                            @php $paragraph[] = $verse; @endphp
                        @endforeach
                        @if ($paragraph !== [])
                            <p>
                                @foreach ($paragraph as $index => $part)
                                    @if ($index > 0)&nbsp;@endif
                                    <sup>{{ $part['number'] }}</sup>{{ $part['text'] }}
                                @endforeach
                            </p>
                        @endif
                    </div>
                @elseif (($column['kind'] ?? '') === 'notes')
                    <div class="notes-content">{{ $column['content'] }}</div>
                @elseif (($column['kind'] ?? '') === 'lined-notes')
                    <div class="column-subtitle">Notes</div>
                    <div class="lined-block" aria-hidden="true"></div>
                @elseif (($column['kind'] ?? '') === 'scribe')
                    <div class="column-subtitle">Scribe</div>
                    @foreach ($column['linedVerses'] as $verse)
                        <div @class(['scribe-verse', 'paragraph-start' => $verse['paragraphStart'] ?? false])>
                            <sup>{{ $verse['number'] }}</sup>
                            <span class="line"></span>
                        </div>
                    @endforeach
                @elseif (($column['kind'] ?? '') === 'message')
                    <p class="message">{{ $column['message'] }}</p>
                @endif
            </div>
        @endforeach
    </div>
</body>
</html>
