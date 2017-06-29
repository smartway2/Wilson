let initialState = []

export default (state = initialState, payload) => {
    switch (payload.type) {
        case 'add':
            return [
                ...state,
                {
                    title: payload.item,
                    completed: false,
                    id: makeId()
                }
            ];
        case 'remove':
            return state.filter(item => item !== payload.item)
        case 'complete':
            return [...state]
        default:
            return state;
    }
};

var counter = 0
function makeId() {
    return counter++;
}
