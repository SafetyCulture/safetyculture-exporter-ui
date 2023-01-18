export function trim(value, size) {
    if (value.length > size) {
        return value.substring(0, size).concat(" ...")
    }
    return value
}
