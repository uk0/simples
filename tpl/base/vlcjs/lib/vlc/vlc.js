import { MediaPlayer } from "./libvlc.js";

export async function VLCPlayer(video,canvas_Id,overlay_id='overlay') {
    let canvas = document.getElementById(canvas_Id);
    console.log(canvas.id);
    if (canvas === null) {
        console.error("No canvas with id 'canvas' found.");
        return;
    }

    let vlcModule = await initModule({
        vlc_access_file: {},
        canvas: canvas,
        canvas_id: canvas_Id,
        overlay_id: overlay_id,
    })

    let media_player;

    let vlc_opts_array = video.options.split(' ');

    let vlc_opts_size = 0;
    for (let i in vlc_opts_array) {
        vlc_opts_size += vlc_opts_array[i].length + 1;
    }

    let buffer = vlcModule._malloc(vlc_opts_size);
    let wrote_size = 0;
    for (let i in vlc_opts_array) {
        vlcModule.writeAsciiToMemory(vlc_opts_array[i], buffer + wrote_size, false);
        wrote_size += vlc_opts_array[i].length + 1;
    }

    let vlc_argv = vlcModule._malloc(vlc_opts_array.length * 4 + 4);
    let view_vlc_argv = new Uint32Array(
        vlcModule.wasmMemory.buffer,
        vlc_argv,
        vlc_opts_array.length
    );

    wrote_size = 0;
    for (let i in vlc_opts_array) {
        view_vlc_argv[i] = buffer + wrote_size;
        wrote_size += vlc_opts_array[i].length + 1;
    }
    console.log(" vlc_argv ",vlc_argv)
    vlcModule._wasm_libvlc_init(vlc_opts_array.length, vlc_argv);
    media_player = new MediaPlayer(vlcModule, "emjsfile://1");
    vlcModule._set_global_media_player(media_player.media_player_ptr);
    media_player.vlcModule = vlcModule

    return media_player
}

