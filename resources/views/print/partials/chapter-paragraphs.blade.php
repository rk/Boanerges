@php
    $paragraph = [];
@endphp
@foreach ($verses as $verse)
    @if (($verse['paragraphStart'] ?? false) && $paragraph !== [])
        <p>
            @foreach ($paragraph as $index => $part)
                @if ($index > 0)&nbsp;@endif
                <sup>{{ $part['number'] }}</sup>@if ($escapeText ?? false){{ $part['text'] }}@else{!! $part['text'] !!}@endif
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
            <sup>{{ $part['number'] }}</sup>@if ($escapeText ?? false){{ $part['text'] }}@else{!! $part['text'] !!}@endif
        @endforeach
    </p>
@endif
