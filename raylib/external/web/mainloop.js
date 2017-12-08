function _emscripten_set_main_loop_go(fps, simulateInfiniteLoop) {
    if (_go_update_function !== undefined && typeof _go_update_function == 'function') {
        _emscripten_set_main_loop(_go_update_function, fps, simulateInfiniteLoop);
    }
}

Module["_emscripten_set_main_loop_go"] = _emscripten_set_main_loop_go;
